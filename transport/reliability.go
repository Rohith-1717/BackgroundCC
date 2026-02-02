package transport

import (
	"backgroundcc/ledbatpp"
	"net"
	"sync"
	"time"
)

type SentPacket struct {
	Packet *Packet
	Addr   *net.UDPAddr
	SentAt time.Time
}

type ReliableTransport struct {
	UDP      *UDPTransport
	Sampler  *ledbatpp.Sampler
	Tracker  *ledbatpp.DelayTracker
	Loss     *ledbatpp.Loss
	State    *ledbatpp.State
	Mutex    sync.Mutex
	NextSeq  uint64
	Unacked  map[uint64]*SentPacket
	Timeout  time.Duration
	StopChan chan struct{}
}

func NewReliableTransport(
	udp *UDPTransport,
	sampler *ledbatpp.Sampler,
	tracker *ledbatpp.DelayTracker,
	loss *ledbatpp.Loss,
	state *ledbatpp.State,
) *ReliableTransport {
	rt := &ReliableTransport{
		UDP:      udp,
		Sampler:  sampler,
		Tracker:  tracker,
		Loss:     loss,
		State:    state,
		NextSeq:  1,
		Unacked:  make(map[uint64]*SentPacket),
		Timeout:  500 * time.Millisecond,
		StopChan: make(chan struct{}),
	}
	go rt.RetransmitLoop()
	return rt
}

func (rt *ReliableTransport) Send(addr *net.UDPAddr, payload []byte) error {
	rt.Mutex.Lock()
	seq := rt.NextSeq
	rt.NextSeq++
	pkt := &Packet{
		Type:      PacketTypeData,
		Seq:       seq,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	rt.Unacked[seq] = &SentPacket{
		Packet: pkt,
		Addr:   addr,
		SentAt: time.Now(),
	}
	rt.Mutex.Unlock()
	return rt.UDP.Send(addr, pkt)
}

func (rt *ReliableTransport) HandleIncoming(pkt *Packet, addr *net.UDPAddr) {
	switch pkt.Type {
	case PacketTypeAck:
		rt.handleAck(pkt)
	case PacketTypeData:
		rt.sendAck(pkt.Seq, pkt.Timestamp, addr)
	}
}

func (rt *ReliableTransport) handleAck(pkt *Packet) {
	var sendTime time.Time
	var ok bool
	rt.Mutex.Lock()
	_, ok = rt.Unacked[pkt.Ack]
	if ok {
		delete(rt.Unacked, pkt.Ack)
	}
	rt.Mutex.Unlock()
	if !ok {
		return
	}
	sendTime = time.Unix(0, pkt.Timestamp)
	sample, ok := rt.Sampler.Observe(sendTime)
	if !ok {
		return
	}
	rt.Tracker.OnRTTSample(sample.RTT)
}

func (rt *ReliableTransport) sendAck(seq uint64, ts int64, addr *net.UDPAddr) {
	ack := &Packet{
		Type:      PacketTypeAck,
		Ack:       seq,
		Timestamp: ts,
	}
	_ = rt.UDP.Send(addr, ack)
}

func (rt *ReliableTransport) RetransmitLoop() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			rt.Mutex.Lock()
			for _, sp := range rt.Unacked {
				if now.Sub(sp.SentAt) >= rt.Timeout {
					rt.Loss.OnLoss(rt.State)
					sp.SentAt = now
					_ = rt.UDP.Send(sp.Addr, sp.Packet)
				}
			}
			rt.Mutex.Unlock()
		case <-rt.StopChan:
			return
		}
	}
}

func (rt *ReliableTransport) Close() {
	close(rt.StopChan)
	_ = rt.UDP.Close()
}

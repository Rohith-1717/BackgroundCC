package transport

import (
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
	Mutex    sync.Mutex
	NextSeq  uint64
	Unacked  map[uint64]*SentPacket
	Timeout  time.Duration
	StopChan chan struct{}
}

func NewReliableTransport(UDP *UDPTransport) *ReliableTransport {
	RT := &ReliableTransport{
		UDP:      UDP,
		NextSeq:  1,
		Unacked:  make(map[uint64]*SentPacket),
		Timeout:  500 * time.Millisecond,
		StopChan: make(chan struct{}),
	}
	go RT.RetransmitLoop()
	return RT
}

func (RT *ReliableTransport) Send(Addr *net.UDPAddr, Payload []byte) error {
	RT.Mutex.Lock()
	Seq := RT.NextSeq
	RT.NextSeq++
	Pkt := &Packet{
		Type:      PacketTypeData,
		Seq:       Seq,
		Ack:       0,
		Timestamp: time.Now().UnixNano(),
		Payload:   Payload,
	}

	RT.Unacked[Seq] = &SentPacket{
		Packet: Pkt,
		Addr:   Addr,
		SentAt: time.Now(),
	}
	RT.Mutex.Unlock()

	return RT.UDP.Send(Addr, Pkt)
}

func (RT *ReliableTransport) HandleIncoming(Pkt *Packet, Addr *net.UDPAddr) {
	switch Pkt.Type {
	case PacketTypeAck:
		RT.HandleAck(Pkt.Ack)
	case PacketTypeData:
		RT.SendAck(Pkt.Seq, Addr)
	}
}

func (RT *ReliableTransport) HandleAck(Ack uint64) {
	RT.Mutex.Lock()
	delete(RT.Unacked, Ack)
	RT.Mutex.Unlock()
}

func (RT *ReliableTransport) SendAck(Seq uint64, Addr *net.UDPAddr) {
	AckPkt := &Packet{
		Type:      PacketTypeAck,
		Seq:       0,
		Ack:       Seq,
		Timestamp: time.Now().UnixNano(),
		Payload:   nil,
	}
	_ = RT.UDP.Send(Addr, AckPkt)
}

func (RT *ReliableTransport) RetransmitLoop() {
	Ticker := time.NewTicker(50 * time.Millisecond)
	defer Ticker.Stop()

	for {
		select {
		case <-Ticker.C:
			Now := time.Now()
			RT.Mutex.Lock()
			for _, SP := range RT.Unacked {
				if Now.Sub(SP.SentAt) >= RT.Timeout {
					SP.SentAt = Now
					_ = RT.UDP.Send(SP.Addr, SP.Packet)
				}
			}
			RT.Mutex.Unlock()
		case <-RT.StopChan:
			return
		}
	}
}

func (RT *ReliableTransport) Close() error {
	close(RT.StopChan)
	return RT.UDP.Close()
}

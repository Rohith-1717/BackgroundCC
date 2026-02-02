package session

import (
	"backgroundcc/ledbatpp"
	"backgroundcc/pacing"
	"backgroundcc/transport"
	"net"
	"sync"
	"time"
)

type SessionState struct {
	Mutex      sync.Mutex
	LocalAddr  *net.UDPAddr
	RemoteAdr  *net.UDPAddr
	Transport  *transport.ReliableTransport
	Pacer      *pacing.Pacer
	Controller *ledbatpp.Controller
	Startup    *ledbatpp.Startup
	Slowdown   *ledbatpp.Slowdown
	Params     ledbatpp.Params
	State      *ledbatpp.State
	StartTime  time.Time
	Closed     bool
	StopChan   chan struct{}
}

func NewSessionState(
	localAddr *net.UDPAddr,
	remoteAddr *net.UDPAddr,
	transport *transport.ReliableTransport,
	pacer *pacing.Pacer,
	controller *ledbatpp.Controller,
	startup *ledbatpp.Startup,
	slowdown *ledbatpp.Slowdown,
	params ledbatpp.Params,
	state *ledbatpp.State,
) *SessionState {

	s := &SessionState{
		LocalAddr:  localAddr,
		RemoteAdr:  remoteAddr,
		Transport:  transport,
		Pacer:      pacer,
		Controller: controller,
		Startup:    startup,
		Slowdown:   slowdown,
		Params:     params,
		State:      state,
		StartTime:  time.Now(),
		Closed:     false,
		StopChan:   make(chan struct{}),
	}
	go s.Transport.ReceiveLoop()
	return s
}

func (s *SessionState) Close() {
	s.Mutex.Lock()
	if s.Closed {
		s.Mutex.Unlock()
		return
	}
	s.Closed = true
	close(s.StopChan)
	s.Transport.Close()
	s.Mutex.Unlock()
}

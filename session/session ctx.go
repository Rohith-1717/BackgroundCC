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
	StartTime  time.Time
	Closed     bool

	StopChan chan struct{}
}

func NewSessionState(
	LocalAddr *net.UDPAddr,
	RemoteAdr *net.UDPAddr,
	Transport *transport.ReliableTransport,
	Pacer *pacing.Pacer,
	Controller *ledbatpp.Controller,
) *SessionState {
	return &SessionState{
		LocalAddr:  LocalAddr,
		RemoteAdr:  RemoteAdr,
		Transport:  Transport,
		Pacer:      Pacer,
		Controller: Controller,
		StartTime:  time.Now(),
		Closed:     false,
		StopChan:   make(chan struct{}),
	}
}

func (S *SessionState) IsClosed() bool {
	S.Mutex.Lock()
	defer S.Mutex.Unlock()
	return S.Closed
}

func (S *SessionState) Close() {
	S.Mutex.Lock()
	if S.Closed {
		S.Mutex.Unlock()
		return
	}
	S.Closed = true
	close(S.StopChan)
	S.Mutex.Unlock()
}

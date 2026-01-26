package session

import (
	"backgroundcc/ledbatpp"
	"io"
	"time"
)

type Sender struct {
	Session *SessionState
	State   *ledbatpp.State
	Reader  io.Reader
	PktSize int
}

func NewSender(
	session *SessionState,
	state *ledbatpp.State,
	reader io.Reader,
	pktSize int,
) *Sender {
	return &Sender{
		Session: session,
		State:   state,
		Reader:  reader,
		PktSize: pktSize,
	}
}

func (s *Sender) Run() error {
	buf := make([]byte, s.PktSize)

	for {
		select {
		case <-s.Session.StopChan:
			return nil
		default:
		}
		now := s.Session.Pacer.Clk.Now()

		// So here we are asking the algo to look at the new network delay state
		// and it will update the sending rate

		s.Session.Controller.Update(s.State)
		s.Session.Pacer.UpdateRate(s.State.Rate)
		if !s.Session.Pacer.CanSend(now, s.PktSize) {
			next := s.Session.Pacer.NextSendTime(s.PktSize)
			if !next.IsZero() {
				sleep := next.Sub(now)
				if sleep > 0 {
					time.Sleep(sleep)
				}
			}
			continue
		}

		n, err := s.Reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				s.Session.Close()
				return nil
			}
			return err
		}
		payload := make([]byte, n)
		copy(payload, buf[:n])
		err = s.Session.Transport.Send(s.Session.RemoteAdr, payload)
		if err != nil {
			return err
		}
		s.Session.Pacer.OnSend(now, n)
	}
}

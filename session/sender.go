package session

import (
	"io"
	"time"
)

type Sender struct {
	Session *SessionState
	Reader  io.Reader
	PktSize int
}

func NewSender(
	session *SessionState,
	reader io.Reader,
	pktSize int,
) *Sender {
	return &Sender{
		Session: session,
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

		s.Session.Startup.Update(
			s.Session.State,
			s.Session.Params,
			now,
		)
		s.Session.Slowdown.MaybeEnter(s.Session.State)
		s.Session.Controller.Update(s.Session.State)
		s.Session.Slowdown.Apply(s.Session.State)
		s.Session.Pacer.UpdateRate(s.Session.State.Rate)

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
		if err := s.Session.Transport.Send(s.Session.RemoteAdr, payload); err != nil {
			return err
		}
		s.Session.Pacer.OnSend(now, n)
	}
}

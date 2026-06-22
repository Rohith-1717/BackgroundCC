package transport

import (
	"os"
	"sync"
)

type Receiver struct {
	mu          sync.Mutex
	file        *os.File
	expectedSeq uint64
	pending     map[uint64][]byte
	done        bool
}

func NewReceiver(path string) (*Receiver, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &Receiver{
		file:        f,
		expectedSeq: 1,
		pending:     make(map[uint64][]byte),
		done:        false,
	}, nil
}

func (r *Receiver) HandleData(pkt *Packet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.done {
		return nil
	}

	if pkt.Seq < r.expectedSeq {
		return nil
	}

	if pkt.Seq > r.expectedSeq {
		if _, exists := r.pending[pkt.Seq]; !exists {
			payload := make([]byte, len(pkt.Payload))
			copy(payload, pkt.Payload)
			r.pending[pkt.Seq] = payload
		}
		return nil
	}

	if _, err := r.file.Write(pkt.Payload); err != nil {
		return err
	}

	r.expectedSeq++

	for {
		payload, exists := r.pending[r.expectedSeq]
		if !exists {
			break
		}
		if _, err := r.file.Write(payload); err != nil {
			return err
		}
		delete(r.pending, r.expectedSeq)
		r.expectedSeq++
	}

	return nil
}

func (r *Receiver) HandleEOF() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.done {
		return nil
	}

	r.done = true
	return r.file.Close()
}

func (r *Receiver) Done() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.done
}

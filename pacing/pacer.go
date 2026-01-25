package pacing

import(
	"math"
	"sync"
	"time"
)

type Clock interface{
	Now() time.Time
}

type Pacer struct{
	Mu sync.Mutex
	Clk Clock
	BytesPerSec float64
	LastSend time.Time
	HasSent bool
}

func NewPacer(clk Clock) *Pacer{
	return &Pacer{
		Clk: clk,
		BytesPerSec: 0,
		HasSent: false,
	}
}

func (p *Pacer) UpdateRate(BytesPerSec float64){
	p.Mu.Lock()
	defer p.Mu.Unlock()
	if BytesPerSec < 0{
		BytesPerSec = 0
	}
	p.BytesPerSec = BytesPerSec
}

func (p *Pacer) interval(pktSize int) time.Duration{
	if (p.BytesPerSec <= 0 || pktSize <= 0){
		return time.Duration(math.MaxInt64)
	}
	seconds := float64(pktSize)/p.BytesPerSec
	return time.Duration(seconds*float64(time.Second))
}

func (p *Pacer) CanSend(now time.Time, pktSize int) bool{
	p.Mu.Lock()
	defer p.Mu.Unlock()
	if p.BytesPerSec <= 0{
		return false
	}
	if !p.HasSent{
		return true
	}
	next := p.LastSend.Add(p.interval(pktSize))
	return !now.Before(next)
}

func (p *Pacer) OnSend(now time.Time, pktSize int){
	p.Mu.Lock()
	defer p.Mu.Unlock()
	if !p.HasSent{
		p.LastSend = now
		p.HasSent = true
		return
	}

	iv := p.interval(pktSize)
	if (iv <= 0){
		p.LastSend = now
		return
	}

	expected := p.LastSend.Add(iv)
	if now.After(expected){
		p.LastSend = now
	} else{
		p.LastSend = expected
	}
}

func (p *Pacer) NextSendTime(pktSize int) time.Time{
	p.Mu.Lock()
	defer p.Mu.Unlock()
	now := p.Clk.Now()
	if p.BytesPerSec <= 0{
		return time.Time{}
	}
	if !p.HasSent{
		return now
	}

	return p.LastSend.Add(p.interval(pktSize))
}

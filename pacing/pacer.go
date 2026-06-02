package pacing
import(
	"math"
	"sync"
	"time"
)

// The pacer decides when the next packet has to be sent after getting the rate from ledbatpp

type Clock interface{
	Now() time.Time
}

type Pacer struct{
	Mu sync.Mutex  // stops concurrent access
	Clk Clock
	BytesPerSec float64  // sending rate in bytes per sec
	LastSend time.Time
	HasSent bool // whether even one packet has been sent till now
}

func NewPacer(clk Clock) *Pacer{  // constructor
	return &Pacer{
		Clk: clk,
		BytesPerSec: 0,
		HasSent: false,
	}
}

func (p *Pacer) UpdateRate(BytesPerSec float64){
	p.Mu.Lock()  // lock it
	defer p.Mu.Unlock()  // unlock once functuon ends
	if BytesPerSec < 0{
		BytesPerSec = 0
	}
	p.BytesPerSec = BytesPerSec
}

func (p *Pacer) interval(pktSize int) time.Duration{  // calculates the interval, i.e. if the current rate is x bytes per sec, how long can I wait between packets
	if (p.BytesPerSec <= 0 || pktSize <= 0){
		return time.Duration(math.MaxInt64)  // i.e. never send
	}
	seconds := float64(pktSize)/p.BytesPerSec  // time needed to send one packet at current rate
	return time.Duration(seconds*float64(time.Second))
}

func (p *Pacer) CanSend(now time.Time, pktSize int) bool{  // can I send it now if I tried to?
	p.Mu.Lock()  // lock
	defer p.Mu.Unlock()  // unlock when function exits
	if p.BytesPerSec <= 0{
		return false
	}
	if !p.HasSent{  // If I've never sent before then send
		return true
	}
	next := p.LastSend.Add(p.interval(pktSize))  // computes next allowed send time
	return !now.Before(next)  // send if its after the next allowed end time
}

func (p *Pacer) OnSend(now time.Time, pktSize int){  // update pacing schedule once that I've sent
	p.Mu.Lock()
	defer p.Mu.Unlock()
	if !p.HasSent{  // if its the first packet, then update
		p.LastSend = now  
		p.HasSent = true
		return
	}

	iv := p.interval(pktSize)  // expected spacing between packets
	
	if (iv <= 0){  // Invalid interval, use current time
		p.LastSend = now
		return
	}

	expected := p.LastSend.Add(iv)  // expected next send time, its expected cuz there might be delays
	if now.After(expected){ // if now is after the expected (maybe cuz of OS delays)
		p.LastSend = now  // last sent time is now
	} else{
		p.LastSend = expected
	}
}

func (p *Pacer) NextSendTime(pktSize int) time.Time{  // gives when I can send
	p.Mu.Lock()
	defer p.Mu.Unlock()
	now := p.Clk.Now()
	if p.BytesPerSec <= 0{
		return time.Time{}
	}
	if !p.HasSent{  // u can send now
		return now
	}

	return p.LastSend.Add(p.interval(pktSize))  // last time a packet was sent + interval
}

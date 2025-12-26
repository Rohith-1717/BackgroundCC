package ledbatpp

import (
	"time"
	"backgroundcc/internal/pacing"
)

type Sample struct{
	SendTime time.Time
	AckTime time.Time
	RTT time.Duration
}

type Sampler struct{
	clock pacing.Clock
}

func NewSampler(clock pacing.Clock) *Sampler{
	return &Sampler{
		clock: clock,
	}
}

func (s *Sampler) Observe(sendTime time.Time)(Sample, bool){
	now := s.clock.Now()

	if sendTime.After(now){
		return Sample{}, false
	}	
	rtt := now.Sub(sendTime)
	if rtt <= 0{
		return Sample{}, false
	}
	return Sample{SendTime: sendTime, AckTime: now,RTT: rtt,}, true
}
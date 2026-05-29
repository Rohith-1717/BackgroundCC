package ledbatpp
import "time"

// This computes the rtt

type Sample struct{
	SendTime time.Time  // time we sent
	AckTime time.Time  // time we got ACK
	RTT time.Duration  // ack time - send time
}

type Sampler struct{
	clock Clock  // Sampler's clock
}

func NewSampler(clock Clock) *Sampler{
	return &Sampler{clock: clock,}
}

func (s *Sampler) Observe(sendTime time.Time)(Sample, bool){
	now := s.clock.Now()  // ack arrival time 

	if sendTime.After(now){
		return Sample{}, false  // if sent time is after now, then its not possible since sent time is after ack recieved time
	}

	rtt := now.Sub(sendTime)
	if rtt <= 0{
		return Sample{}, false
	}

	return Sample{SendTime: sendTime,	AckTime: now, RTT: rtt,}, true  // return the sent time, ack time and rtt
}

// This is to make sure that the RTTs that we observed make sense

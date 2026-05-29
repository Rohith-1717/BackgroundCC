package ledbatpp
import "time"

// Instead of storing just one base delay, we have 2 windows long term and short term, with rtt samples
// We get the base delay from the long term window, i.e. minimum rtt in long term window
// We get current delay from the short term window, i.e. minimum rtt in the short term window
// Queue delay = current delay - base delay 

type delaySample struct{  // one rtt measurement sample
	ts  time.Time
	rtt time.Duration
}

type DelayEstimator struct{
	baseWindow time.Duration  // long term history window
	currWindow time.Duration  // short term history window
	baseSamples []delaySample
	currSamples []delaySample
}

func NewDelayEstimator(baseWindow, currWindow time.Duration) *DelayEstimator{  // constructor
	return &DelayEstimator{
		baseWindow: baseWindow,
		currWindow: currWindow,
	}
}

// Add new sample to window
func (d *DelayEstimator) Update(rtt time.Duration, now time.Time){
	d.baseSamples = append(d.baseSamples, delaySample{ts:  now, rtt: rtt,})
	d.currSamples = append(d.currSamples, delaySample{ts:  now,	rtt: rtt,})
	d.expireOld(now)  // remove old sample 
}

func (d *DelayEstimator) expireOld(now time.Time){
	// window cleanup 
	baseCutoff := now.Add(-d.baseWindow)  // compute the oldest allowed timestamp
	currCutoff := now.Add(-d.currWindow)
	i := 0
	for ; i < len(d.baseSamples); i++{
		if d.baseSamples[i].ts.After(baseCutoff){  // find the first non-expired sample
			break
		}
	}
	d.baseSamples = d.baseSamples[i:]  // the base samples start from i
	i = 0
	for ; i < len(d.currSamples); i++{
		if d.currSamples[i].ts.After(currCutoff){  // same thing with short term window
			break
		}
	}
	d.currSamples = d.currSamples[i:]
}

func minRTT(samples []delaySample) time.Duration{
	if len(samples) == 0 {
		return 0
	}
	min := samples[0].rtt  // initialize min rtt
	for _, s := range samples[1:]{
		if s.rtt < min{
			min = s.rtt  // find min rtt
		}
	}
	return min  // return min rtt
}

func (d *DelayEstimator) BaseDelay() time.Duration{  // return the base delay
	return minRTT(d.baseSamples)
}

func (d *DelayEstimator) CurrentDelay() time.Duration{  // return the current delay
	return minRTT(d.currSamples)
}

func (d *DelayEstimator) QueuingDelay() time.Duration{
	base := d.BaseDelay()
	curr := d.CurrentDelay()

	if base == 0 || curr <= base{
		return 0
	}
	return curr - base  // This is the queue delay
}

package ledbatpp
import "time"

type delaySample struct{
	ts  time.Time
	rtt time.Duration
}

type DelayEstimator struct{
	baseWindow time.Duration
	currWindow time.Duration
	baseSamples []delaySample
	currSamples []delaySample
}

func NewDelayEstimator(baseWindow, currWindow time.Duration) *DelayEstimator{
	return &DelayEstimator{
		baseWindow: baseWindow,
		currWindow: currWindow,
	}
}

func (d *DelayEstimator) Update(rtt time.Duration, now time.Time){
	d.baseSamples = append(d.baseSamples, delaySample{ts:  now, rtt: rtt,})
	d.currSamples = append(d.currSamples, delaySample{ts:  now,	rtt: rtt,})
	d.expireOld(now)
}

func (d *DelayEstimator) expireOld(now time.Time){
	baseCutoff := now.Add(-d.baseWindow)
	currCutoff := now.Add(-d.currWindow)
	i := 0
	for ; i < len(d.baseSamples); i++{
		if d.baseSamples[i].ts.After(baseCutoff){
			break
		}
	}
	d.baseSamples = d.baseSamples[i:]
	i = 0
	for ; i < len(d.currSamples); i++{
		if d.currSamples[i].ts.After(currCutoff){
			break
		}
	}
	d.currSamples = d.currSamples[i:]
}

func minRTT(samples []delaySample) time.Duration{
	if len(samples) == 0 {
		return 0
	}
	min := samples[0].rtt
	for _, s := range samples[1:]{
		if s.rtt < min{
			min = s.rtt
		}
	}
	return min
}

func (d *DelayEstimator) BaseDelay() time.Duration{
	return minRTT(d.baseSamples)
}

func (d *DelayEstimator) CurrentDelay() time.Duration{
	return minRTT(d.currSamples)
}

func (d *DelayEstimator) QueuingDelay() time.Duration{
	base := d.BaseDelay()
	curr := d.CurrentDelay()

	if base == 0 || curr <= base{
		return 0
	}
	return curr - base
}

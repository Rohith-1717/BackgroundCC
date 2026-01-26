package ledbatpp

import "time"

type DelayTracker struct {
	estimator *DelayEstimator
	state     *State
	clock     Clock
}

func NewDelayTracker(
	estimator *DelayEstimator,
	state *State,
	clock Clock,
) *DelayTracker {
	return &DelayTracker{
		estimator: estimator,
		state:     state,
		clock:     clock,
	}
}

func (d *DelayTracker) OnRTTSample(rtt time.Duration) {
	now := d.clock.Now()
	d.estimator.Update(rtt, now)
	d.state.BaseDelay = d.estimator.BaseDelay()
	d.state.CurrentDelay = d.estimator.CurrentDelay()
	d.state.QueuingDelay = d.estimator.QueuingDelay()
}

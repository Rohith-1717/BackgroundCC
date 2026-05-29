package ledbatpp
import "time"

// This updates the state by taking data from the delay estimator and giving it to the state so the controller is updated

type DelayTracker struct {
	estimator *DelayEstimator
	state *State
	clock Clock
}

func NewDelayTracker(
	estimator *DelayEstimator,
	state *State,
	clock Clock,
) *DelayTracker {
	return &DelayTracker{
		estimator: estimator,
		state: state,
		clock: clock,
	}
}

func (d *DelayTracker) OnRTTSample(rtt time.Duration) {
	now := d.clock.Now()  // get the current time
	d.estimator.Update(rtt, now)  // update the delay estimator
	d.state.BaseDelay = d.estimator.BaseDelay()  // update the values onto state
	d.state.CurrentDelay = d.estimator.CurrentDelay()
	d.state.QueuingDelay = d.estimator.QueuingDelay()
}

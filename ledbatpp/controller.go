package ledbatpp

import "math"

type Controller struct {
	params Params  // controller parameters
	clock  Clock
}

func NewController(params Params, clock Clock) *Controller {  // constructor
	return &Controller{
		params: params,
		clock:  clock,
	}
}

func (c *Controller) Update(state *State) {
	now := c.clock.Now()  // current time 
	dt := now.Sub(state.LastRateUpdate).Seconds()  // delta t = dt = now - lastupdate
	if dt <= 0 {
		return
	}

	// This is the ledbat error signal
	// This one tells us if we should increase or decrease the sending rate
	// So its the most important part xD

	err := state.TargetDelay - state.QueuingDelay   // target delay - queuing delay
	norm := err.Seconds() / state.TargetDelay.Seconds()  // normalize it for controller independent scaling
	gain := math.Abs(norm)
	if err > 0 {
		state.Rate += c.params.AdditiveIncrease * gain * dt  // if positive error, i.e. queuing delay is less than target delay, increase rate of sending
	} else {
		state.Rate -= state.Rate * c.params.ProportionalDecrease * gain * dt
	}
	if state.Rate < c.params.MinRate {  // maximum rate clamp, i.e. if rate is less than minimum, then clamp it
		state.Rate = c.params.MinRate  // rate is clamped
	}
	if state.Rate > c.params.MaxRate {  // same if it exceeds
		state.Rate = c.params.MaxRate
	}
	state.LastRateUpdate = now  // the last it was updated is now (needed for next dt computation)
}

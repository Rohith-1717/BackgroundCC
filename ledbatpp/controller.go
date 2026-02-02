package ledbatpp

import "math"

type Controller struct {
	params Params
	clock  Clock
}

func NewController(params Params, clock Clock) *Controller {
	return &Controller{
		params: params,
		clock:  clock,
	}
}

func (c *Controller) Update(state *State) {
	now := c.clock.Now()
	dt := now.Sub(state.LastRateUpdate).Seconds()
	if dt <= 0 {
		return
	}

	// This is the ledbat error signal
	// This one tells us if we should increase or decrease the sending rate
	// So its the most important part xD

	err := state.TargetDelay - state.QueuingDelay
	norm := err.Seconds() / state.TargetDelay.Seconds()
	gain := math.Abs(norm)
	if err > 0 {
		state.Rate += c.params.AdditiveIncrease * gain * dt
	} else {
		state.Rate -= state.Rate * c.params.ProportionalDecrease * gain * dt
	}
	if state.Rate < c.params.MinRate {
		state.Rate = c.params.MinRate
	}
	if state.Rate > c.params.MaxRate {
		state.Rate = c.params.MaxRate
	}
	state.LastRateUpdate = now
}

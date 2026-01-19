package ledbatpp
import "math"

type Controller struct{
	params Params
	clock Clock
}

func NewController(params Params, clock Clock) *Controller{
	return &Controller{
		params: params,
		clock: clock,
	}
}

func (c *Controller) Update(state *State){
	now := c.clock.Now()
	dt := now.Sub(state.LastUpdate).Seconds()
	if dt <= 0{
		return
	}

	// This is the ledbat error signal
	// This one tells us if we should increase or decrease the sending rate
	// So its the most important part xD

	error := state.TargetDelay - state.QueuingDelay
	normalizedError := error.Seconds()/state.TargetDelay.Seconds()

	// This is non linear gain scaling

	gain := math.Abs(normalizedError)

	var rateDelta float64

	if error > 0{
		
		rateDelta = c.params.AdditiveIncrease*gain
		state.Rate += rateDelta*dt
	} else{
		rateDelta = c.params.ProportionalDecrease*gain
		state.Rate -= state.Rate*rateDelta*dt
	}
	
	// This part is the clamp rate 
	// this thing forces the sending rate to stay within the correct allowed bounds

	if state.Rate < c.params.MinRate{
		state.Rate = c.params.MinRate
	}
	if state.Rate > c.params.MaxRate{
		state.Rate = c.params.MaxRate
	}

	state.LastUpdate = now
}

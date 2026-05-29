package ledbatpp

import "time"

type Slowdown struct {
	params Params
	clock Clock
	interval time.Duration
	duration time.Duration
}

func NewSlowdown(params Params, clock Clock) *Slowdown {
	return &Slowdown{
		params: params,
		clock: clock,
		interval: 5*time.Second,
		duration: 200*time.Millisecond,
	}
}

func (s *Slowdown) MaybeEnter(state *State) {  // Should we slowdown?
	now := s.clock.Now()
	if state.InStartup || state.InSlowdown {  // dont slow down if we are in startup or already in slowdown
		return
	}
	if state.QueuingDelay < time.Duration(float64(state.TargetDelay)*s.params.SlowdownEnterThreshold) {  // dont go to slowdown if the queue delay is not high enough
		return
	}
	if now.Sub(state.LastSlowdownTime) < s.interval {  // if the last slowdown state wasn't too long ago
		return
	}
	state.InSlowdown = true  // enter slowdown
	state.LastSlowdownTime = now  // last slowdown time is now
}

func (s *Slowdown) Apply(state *State) {
	if !state.InSlowdown {  // If we aren't in slowdown
		return
	}
	state.Rate *= s.params.MultiplicativeDecrease  // do multiplicative decrease
	if state.Rate < s.params.MinRate {
		state.Rate = s.params.MinRate  // clamp the min rate
	}
	if s.clock.Since(state.LastSlowdownTime) >= s.duration {  // If the slowdown has been going on for sometime, then stop
		state.InSlowdown = false
		state.LastSlowdownTime = s.clock.Now()
	}
}

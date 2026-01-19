package ledbatpp
import "time"

type Slowdown struct{
	params Params
	clock Clock
	interval time.Duration
	duration time.Duration
}
func NewSlowdown(params Params, clock Clock) *Slowdown{
	return &Slowdown{
		params: params,
		clock: clock,
		interval: 5*time.Second,
		duration: 200*time.Millisecond,
	}
}

func (s *Slowdown) MaybeEnter(state *State){
	if state.InStartup || state.InSlowdown{
		return
	}
	if state.QueuingDelay < time.Duration(
		float64(state.TargetDelay)*s.params.SlowdownEnterThreshold,
	){
		return
	}
	if s.clock.Since(state.LastUpdate) < s.interval{
		return
	}
	state.InSlowdown = true
	state.LastUpdate = s.clock.Now()
}

func (s *Slowdown) Apply(state *State){
	if !state.InSlowdown{
		return
	}
	state.Rate *= s.params.MultiplicativeDecrease
	if state.Rate < s.params.MinRate{
		state.Rate = s.params.MinRate
	}
	if s.clock.Since(state.LastUpdate) >= s.duration{
		state.InSlowdown = false
		state.LastUpdate = s.clock.Now()
	}
}

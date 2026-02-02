package ledbatpp

import "time"

type Loss struct {
	params   Params
	clock    Clock
	cooldown time.Duration
}

func NewLoss(params Params, clock Clock) *Loss {
	return &Loss{
		params:   params,
		clock:    clock,
		cooldown: 500 * time.Millisecond,
	}
}

func (l *Loss) OnLoss(state *State) {
	now := l.clock.Now()
	if state.InStartup {
		return
	}
	if now.Sub(state.LastLossTime) < l.cooldown {
		return
	}
	if state.QueuingDelay < state.TargetDelay {
		return
	}
	state.Rate *= l.params.MultiplicativeDecrease
	if state.Rate < l.params.MinRate {
		state.Rate = l.params.MinRate
	}
	state.InSlowdown = false
	state.LastLossTime = now
}

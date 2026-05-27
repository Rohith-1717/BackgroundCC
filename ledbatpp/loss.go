package ledbatpp

import "time"

// Essentially this is an emergency secondaruy congestion controller
// Normal TCP works with packet loss
// If packet loss occurs here, loss.go deals with it
// Packet loss still signals severe congestion, queue overflow and unstable links

type Loss struct {
	params   Params
	clock    Clock
	cooldown time.Duration  // This is a cooldown so that it doesn't keep reacting repeatedly to bursts of packet loss
}

func NewLoss(params Params, clock Clock) *Loss {  // constructor
	return &Loss{
		params:   params,
		clock:    clock,
		cooldown: 500 * time.Millisecond,  // Allow atmost one loss reaction every 500 ms
	}
}

func (l *Loss) OnLoss(state *State) {
	now := l.clock.Now() // current time
	if state.InStartup {  // Ignore loss during startup, cuz temporary packet drops may happen naturally
		return
	}
	if now.Sub(state.LastLossTime) < l.cooldown {  // cooldown check
		return
	}
	if state.QueuingDelay < state.TargetDelay {  // it should react only if its congested
		return
	}
	state.Rate *= l.params.MultiplicativeDecrease  // Multiplicative rate decrease if packet loss happens
	if state.Rate < l.params.MinRate {  // clamp the min rate
		state.Rate = l.params.MinRate
	}
	state.InSlowdown = false  // exit slowdown mode because when real congestion happens, we dont have to periodically slowdown
	state.LastLossTime = now
}

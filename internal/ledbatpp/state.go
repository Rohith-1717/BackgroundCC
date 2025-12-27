package ledbatpp

import "time"

type State struct{
	BaseDelay time.Duration
	CurrentDelay time.Duration
	QueuingDelay time.Duration

	TargetDelay time.Duration
	Rate float64
	LastUpdate time.Time
	
	InStartup bool
	InSlowdown bool
}

func NewState(targetDelay time.Duration, initialRate float64, now time.Time) *State{
	return &State{
		BaseDelay: 0,
		CurrentDelay: 0,
		QueuingDelay: 0,
		TargetDelay: targetDelay,
		Rate: initialRate,
		LastUpdate: now,
		InStartup: true,
		InSlowdown: false,
	}
}

// This is basically the memory of LEDBAT++

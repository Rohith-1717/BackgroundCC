package ledbatpp

import "time"

type State struct {
	BaseDelay        time.Duration
	CurrentDelay     time.Duration
	QueuingDelay     time.Duration
	TargetDelay      time.Duration
	Rate             float64
	LastRateUpdate   time.Time
	LastLossTime     time.Time
	LastSlowdownTime time.Time
	InStartup        bool
	InSlowdown       bool
}

func NewState(targetDelay time.Duration, initialRate float64, now time.Time) *State {
	return &State{
		TargetDelay:      targetDelay,
		Rate:             initialRate,
		LastRateUpdate:   now,
		LastLossTime:     now,
		LastSlowdownTime: now,
		InStartup:        true,
		InSlowdown:       false,
	}
}

// This is basically the memory of LEDBAT++

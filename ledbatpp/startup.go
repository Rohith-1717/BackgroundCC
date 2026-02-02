package ledbatpp

import "time"

type Startup struct {
	clock     Clock
	startTime time.Time
	active    bool
}

func NewStartup(clock Clock, now time.Time) *Startup {
	return &Startup{
		clock: clock, startTime: now, active: true,
	}
}

func (s *Startup) Active() bool {
	return s.active
}

func (s *Startup) Update(st *State, p Params, now time.Time) {
	if !s.active {
		return
	}
	if st.TargetDelay <= 0 {
		return
	}
	if st.QueuingDelay > 0 {
		delayRatio := float64(st.QueuingDelay) / float64(st.TargetDelay)
		if delayRatio >= p.StartupExitThreshold {
			s.end(st)
			return
		}
	}
	if s.clock.Since(s.startTime) >= p.CurrentDelayWindow {
		s.end(st)
		return
	}
	increase := p.AdditiveIncrease * p.MinRate
	newRate := st.Rate + increase

	if newRate > p.MaxRate {
		newRate = p.MaxRate
	}

	st.Rate = newRate
}

func (s *Startup) end(st *State) {
	s.active = false
	st.InStartup = false
	st.LastRateUpdate = s.clock.Now()
}

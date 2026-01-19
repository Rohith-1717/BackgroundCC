package ledbatpp
import "time"

type Params struct{
	TargetDelay time.Duration
	BaseDelayWindow time.Duration
	CurrentDelayWindow time.Duration
	AdditiveIncrease float64
	ProportionalDecrease float64
	MultiplicativeDecrease float64
	MinRate float64
	MaxRate float64 
	StartupExitThreshold float64
	SlowdownEnterThreshold float64
}

func DefaultParams() Params{
	return Params{
		TargetDelay: 50*time.Millisecond,
		BaseDelayWindow: 10*time.Minute,
		CurrentDelayWindow: 100*time.Millisecond,
		AdditiveIncrease: 0.05,
		ProportionalDecrease: 0.1,
		MultiplicativeDecrease: 0.7,
		MinRate: 1*1024,
		MaxRate: 10*1024*1024,
		StartupExitThreshold: 0.8,
		SlowdownEnterThreshold: 1.2,
	}
}

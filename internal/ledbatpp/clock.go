package ledbatpp
import "time"

type Clock interface{
	Now() time.Time	
	Since(t time.Time) time.Duration	
	After(d time.Duration) <-chan time.Time	
}

type MonotonicClock struct{}


func NewMonotonicClock() *MonotonicClock{
	return &MonotonicClock{}
}

func (c *MonotonicClock) Now() time.Time{
	return time.Now()
}

func (c *MonotonicClock) Since(t time.Time) time.Duration{
	return time.Since(t)
}

func (c *MonotonicClock) After(d time.Duration) <-chan time.Time{
	return time.After(d)
}


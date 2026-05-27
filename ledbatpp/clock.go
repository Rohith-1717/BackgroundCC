package ledbatpp
import "time"

// This is our clock abstraction

type Clock interface{  // the operations clock must support (IT IS AN INTERFACE BTW)
	Now() time.Time	// current time
	Since(t time.Time) time.Duration // how much time has passed since t
	After(d time.Duration) <-chan time.Time	// for waking up after some delay
}

type MonotonicClock struct{} // this is our real Clock


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


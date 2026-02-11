package main

import (
	"backgroundcc/ledbatpp"
	"backgroundcc/pacing"
	"backgroundcc/session"
	"backgroundcc/transport"
	"bytes"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatalf("usage: %s <local_addr> <remote_addr> <file>", os.Args[0])
	}
	localAddrStr := os.Args[1]
	remoteAddrStr := os.Args[2]
	filePath := os.Args[3]
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	reader := bytes.NewReader(data)
	udp, err := transport.NewUDPTransport(localAddrStr)
	if err != nil {
		log.Fatal(err)
	}
	remoteAddr, err := net.ResolveUDPAddr("udp", remoteAddrStr)
	if err != nil {
		log.Fatal(err)
	}
	clock := ledbatpp.NewMonotonicClock()
	params := ledbatpp.DefaultParams()
	state := ledbatpp.NewState(
		params.TargetDelay,
		params.MinRate,
		clock.Now(),
	)
	controller := ledbatpp.NewController(params, clock)
	startup := ledbatpp.NewStartup(clock, clock.Now())
	slowdown := ledbatpp.NewSlowdown(params, clock)
	sampler := ledbatpp.NewSampler(clock)
	estimator := ledbatpp.NewDelayEstimator(
		params.BaseDelayWindow,
		params.CurrentDelayWindow,
	)
	tracker := ledbatpp.NewDelayTracker(estimator, state, clock)
	loss := ledbatpp.NewLoss(params, clock)
	rt := transport.NewReliableTransport(
		udp,
		sampler,
		tracker,
		loss,
		state,
	)

	pacer := pacing.NewPacer(clock)
	sess := session.NewSessionState(
		nil,
		remoteAddr,
		rt,
		pacer,
		controller,
		startup,
		slowdown,
		params,
		state,
	)
	sender := session.NewSender(sess, reader, 1200)
	log.Println("BackgroundCC started")
	if err := sender.Run(); err != nil {
		log.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	log.Println("BackgroundCC finished")
}

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

func runSender(
	localAddrStr string,
	remoteAddrStr string,
	filePath string,
) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	reader := bytes.NewReader(data)
	udp, err := transport.NewUDPTransport(localAddrStr)
	if err != nil {
		log.Fatal(err)
	}
	remoteAddr, err := net.ResolveUDPAddr(
		"udp",
		remoteAddrStr,
	)
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

	controller := ledbatpp.NewController(
		params,
		clock,
	)
	startup := ledbatpp.NewStartup(
		clock,
		clock.Now(),
	)
	slowdown := ledbatpp.NewSlowdown(
		params,
		clock,
	)
	sampler := ledbatpp.NewSampler(clock)
	estimator := ledbatpp.NewDelayEstimator(
		params.BaseDelayWindow,
		params.CurrentDelayWindow,
	)
	tracker := ledbatpp.NewDelayTracker(
		estimator,
		state,
		clock,
	)
	loss := ledbatpp.NewLoss(
		params,
		clock,
	)
	rt := transport.NewReliableTransport(
		udp,
		sampler,
		tracker,
		loss,
		state,
		nil,
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
	sender := session.NewSender(
		sess,
		reader,
		1200,
	)
	log.Println("BackgroundCC sender started")
	if err := sender.Run(); err != nil {
		log.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	log.Println("BackgroundCC sender finished")
}

func runReceiver(
	localAddrStr string,
	outputPath string,
) {
	receiver, err := transport.NewReceiver(
		outputPath,
	)
	if err != nil {
		log.Fatal(err)
	}
	udp, err := transport.NewUDPTransport(
		localAddrStr,
	)
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

	sampler := ledbatpp.NewSampler(clock)
	estimator := ledbatpp.NewDelayEstimator(
		params.BaseDelayWindow,
		params.CurrentDelayWindow,
	)

	tracker := ledbatpp.NewDelayTracker(
		estimator,
		state,
		clock,
	)

	loss := ledbatpp.NewLoss(
		params,
		clock,
	)

	rt := transport.NewReliableTransport(
		udp,
		sampler,
		tracker,
		loss,
		state,
		receiver,
	)
	go rt.ReceiveLoop()
	log.Println("BackgroundCC receiver started")
	for !receiver.Done() {
		time.Sleep(
			100 * time.Millisecond,
		)
	}
	rt.Close()
	log.Println("BackgroundCC receiver finished")
}

func main() {
	if len(os.Args) >= 2 &&
		os.Args[1] == "validate" {
		runValidate()
		return
	}
	if len(os.Args) < 2 {
		log.Fatal(
			"usage:\n" +
				"backgroundcc send <local_addr> <remote_addr> <file>\n" +
				"backgroundcc receive <local_addr> <output_file>",
		)
	}

	switch os.Args[1] {
	case "send":
		if len(os.Args) != 5 {
			log.Fatal(
				"usage: backgroundcc send <local_addr> <remote_addr> <file>",
			)
		}
		runSender(
			os.Args[2],
			os.Args[3],
			os.Args[4],
		)
	case "receive":
		if len(os.Args) != 4 {
			log.Fatal(
				"usage: backgroundcc receive <local_addr> <output_file>",
			)
		}
		runReceiver(
			os.Args[2],
			os.Args[3],
		) 
	default:
		log.Fatal(
			"unknown mode: use send or receive",
		)
	}
}

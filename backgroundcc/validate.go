package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

func runValidate() {
	fmt.Println("BackgroundCC Validation Report")
	checkSystemInfo()
	checkGo()
	checkGit()
	checkEnv()
	checkDNS()
	checkNetwork()
	checkUDP()
	checkUDPExchange()
	checkPortBinding()
	checkConcurrency()
	fmt.Println("Validation complete.")
}

func checkSystemInfo() {
	fmt.Println("[System]")
	fmt.Printf(" OS : %s\n", runtime.GOOS)
	fmt.Printf(" Arch : %s\n", runtime.GOARCH)
	fmt.Printf(" CPUs : %d\n", runtime.NumCPU())
	fmt.Printf(" Goroutines : %d\n\n", runtime.NumGoroutine())
}

func checkGo() {
	fmt.Println("[Go Environment]")
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		fmt.Println(" FAIL: Go not found\n")
		return
	}
	fmt.Println(" OK  :", strings.TrimSpace(string(out)))
	fmt.Println()
}

func checkGit() {
	fmt.Println("[Git]")
	out, err := exec.Command("git", "--version").Output()
	if err != nil {
		fmt.Println(" WARN: Git not found\n")
		return
	}
	fmt.Println(" OK  :", strings.TrimSpace(string(out)))
	fmt.Println()
}

func checkEnv() {
	fmt.Println("[Environment]")
	path := os.Getenv("PATH")
	if path == "" {
		fmt.Println(" WARN: PATH not set")
	} else {
		fmt.Println(" OK  : PATH detected")
	}
	home := os.Getenv("HOME")
	if home == "" {
		fmt.Println(" WARN: HOME not set")
	} else {
		fmt.Println(" OK  : HOME =", home)
	}
	fmt.Println()
}

func checkDNS() {
	fmt.Println("[DNS Resolution]")
	start := time.Now()
	ips, err := net.LookupHost("google.com")
	if err != nil {
		fmt.Println(" FAIL: DNS resolution failed\n")
		return
	}
	fmt.Printf(" OK  : resolved %d addresses (%v)\n\n", len(ips), time.Since(start))
}

func checkNetwork() {
	fmt.Println("[Network Connectivity]")
	start := time.Now()
	conn, err := net.DialTimeout("tcp", "8.8.8.8:53", 2*time.Second)
	if err != nil {
		fmt.Println(" FAIL: No outbound connectivity\n")
		return
	}
	conn.Close()
	fmt.Printf(" OK  : reachable (latency %v)\n\n", time.Since(start))
}

func checkUDP() {
	fmt.Println("[UDP Capability]")
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(" FAIL: Resolve error\n")
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(" FAIL: Cannot open socket\n")
		return
	}
	defer conn.Close()
	fmt.Println(" OK  : UDP socket created\n")
}

func checkUDPExchange() {
	fmt.Println("[UDP Loopback Test]")
	serverAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println(" FAIL: Cannot start UDP server\n")
		return
	}
	defer serverConn.Close()
	done := make(chan bool)
	go func() {
		buf := make([]byte, 1024)
		n, addr, err := serverConn.ReadFromUDP(buf)
		if err == nil && n > 0 {
			serverConn.WriteToUDP(buf[:n], addr)
			done <- true
		} else {
			done <- false
		}
	}()

	clientConn, err := net.DialUDP("udp", nil, serverConn.LocalAddr().(*net.UDPAddr))
	if err != nil {
		fmt.Println(" FAIL: Client connection failed\n")
		return
	}
	defer clientConn.Close()
	msg := []byte("ping")
	start := time.Now()
	clientConn.Write(msg)
	buf := make([]byte, 1024)
	clientConn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := clientConn.Read(buf)
	if err != nil || n == 0 {
		fmt.Println(" FAIL: No response\n")
		return
	}
	<-done
	fmt.Printf(" OK  : round-trip successful (%v)\n\n", time.Since(start))
}

func checkPortBinding() {
	fmt.Println("[Port Binding]")
	addr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		fmt.Println(" FAIL: Resolve error\n")
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(" FAIL: Bind failed\n")
		return
	}
	defer conn.Close()
	fmt.Println(" OK  : Port binding works\n")
}

func checkConcurrency() {
	fmt.Println("[Concurrency Test]")
	var wg sync.WaitGroup
	start := time.Now()
	count := 100
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			time.Sleep(10 * time.Millisecond)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf(" OK  : handled %d goroutines (%v)\n\n", count, time.Since(start))
}

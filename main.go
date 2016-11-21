// Copyright 2016 Eleme Inc. All rights reserved.

package main

import (
	"fmt"
	logging "log"
	"os"
	"os/signal"
	"syscall"
)

// Variables
var (
	log          = logging.New(os.Stdout, "witch: ", logging.Ldate|logging.Ltime|logging.Lmicroseconds|logging.Lshortfile)
	stopWaitSecs = 5
)

type opts struct {
	confFile string
	execPath string
	execArgs []string
}

func (o *opts) parse() {
	args := os.Args
	if len(args) < 3 {
		log.Printf("More arguments required")
		usage()
		os.Exit(1)
	}
	o.confFile = args[1]
	o.execPath = args[2]
	o.execArgs = args[3:]
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: witch [config file] [cmd path] [cmd arguments]\n")
}

func handleSignals(exitFunc func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-sigs
	log.Printf("Signal %v captured", sig)
	exitFunc()
}

func main() {
	op := &opts{}
	op.parse()

	cfg := &Config{}
	if err := cfg.Parse(op.confFile); err != nil {
		log.Fatalf("Parse config file error: %v", err)
	}

	sys := NewSystem(cfg.PidFile, op.execPath, op.execArgs)
	sys.Start()

	ser := NewServer(cfg.ListenAddr, sys, cfg)
	go func() {
		if err := ser.Start(); err != nil {
			log.Fatalf("Start system server faile: %v", err)
		}
	}()

	handleSignals(func() {
		sys.Stop()
	})
}

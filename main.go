// Copyright 2016 Eleme Inc. All rights reserved.

package main

import (
	"flag"
	logging "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Eagle-X/witch/system"
)

// Variables
var (
	log      = logging.New(os.Stdout, "witch: ", logging.Ldate|logging.Ltime|logging.Lmicroseconds|logging.Lshortfile)
	confFile string
)

func init() {
	flag.StringVar(&confFile, "c", "witch.yaml", "Config file")
	flag.Parse()
}

func handleSignals(exitFunc func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-sigs
	log.Printf("Signal %v captured", sig)
	exitFunc()
}

func createSystem(cfg *Config) system.System {
	switch cfg.Control {
	case "buildin":
		return system.NewLauncher(cfg.PidFile, cfg.Command)
	case "supervisor":
		return system.NewSupervisor(cfg.Service)
	case "systemd":
		return system.NewSystemd(cfg.Service)
	}
	log.Fatalf("Invalid control '%s'", cfg.Control)
	return nil
}

func main() {
	cfg := &Config{}
	if err := cfg.Parse(confFile); err != nil {
		log.Fatalf("Parse config file error: %v", err)
	}

	sys := createSystem(cfg)
	sys.Start()

	ser := NewServer(cfg.ListenAddr, &system.Controller{sys}, cfg)
	go func() {
		if err := ser.Start(); err != nil {
			log.Fatalf("Start system server faile: %v", err)
		}
	}()

	handleSignals(func() {
		sys.Stop()
	})
}

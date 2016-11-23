// Copyright 2016 Eleme Inc. All rights reserved.

package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config is the server configuration.
type Config struct {
	// ListenAddr is the server listen address.
	ListenAddr string `yaml:"listen"`
	// Control is the process control system.
	Control string `yaml:"control"`
	// Service is the name of service controlled by supervisor or systemd.
	Service string `yaml:"service"`
	// Command is the command executed by buildin control.
	Command string `yaml:"command"`
	// PidFile is the process's pid file.
	PidFile string `yaml:"pid_file"`
	// Auth is the username and password configuration.
	Auth map[string]string `yaml:"auth"`
}

// Parse updates config with yaml config file.
func (cfg *Config) Parse(file string) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return err
	}
	return nil
}

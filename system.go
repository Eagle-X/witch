// Copyright 2016 Eleme Inc. All rights reserved.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
	"time"
)

// System supervises the process status, start, stop and restart.
type System struct {
	pidFile  string
	execPath string
	execArgs []string
}

// NewSystem creates new system.
func NewSystem(pidFile, execPath string, execArgs []string) *System {
	return &System{
		pidFile:  pidFile,
		execPath: execPath,
		execArgs: execArgs,
	}
}

func (s *System) writePid(pid int) {
	if err := WriteFile(s.pidFile, []byte(strconv.FormatInt(int64(pid), 10)), 0644); err != nil {
		log.Fatalf("Failed to write pid file: %s", err)
	}
}

func (s *System) readPid() (int, bool) {
	f, err := ioutil.ReadFile(s.pidFile)
	if err != nil {
		log.Printf("Error reading pid file[%s]: %s", s.pidFile, err)
		return -1, false
	}

	pid, err := strconv.Atoi(string(f))
	if err != nil {
		log.Printf("Invalid pid value[%s]: %s", s.pidFile, err)
		return -1, false
	}

	return pid, true
}

func (s *System) pidAlive(pid int) bool {
	return (syscall.Kill(pid, 0) == nil)
}

// IsAlive check if the process alive.
func (s *System) IsAlive() (int, bool) {
	pid, ok := s.readPid()
	if !ok || pid < 1 {
		return pid, false
	}
	return pid, s.pidAlive(pid)
}

// Start starts the process.
func (s *System) Start() (bool, error) {
	if pid, ok := s.IsAlive(); ok {
		log.Printf("The process is alive, pid: %d", pid)
		return true, nil
	}

	log.Printf("Starting %s %s", s.execPath, s.execArgs)
	child := exec.Command(s.execPath, s.execArgs...)
	child.Stdin = os.Stdin
	child.Stdout = os.Stdout
	child.Stderr = os.Stderr
	if err := child.Start(); err != nil {
		log.Printf("Failed to start: %s", err)
		return false, err
	}
	s.writePid(child.Process.Pid)
	go child.Wait()
	return true, nil
}

// Restart restart the process
func (s *System) Restart() (bool, error) {
	s.Stop()
	return s.Start()
}

// Stop stops the process.
func (s *System) Stop() bool {
	pid, ok := s.IsAlive()
	if !ok {
		log.Printf("The process not alive")
		return true
	}
	syscall.Kill(pid, syscall.SIGTERM)
	stopped := make(chan bool)
	go func() {
		for s.pidAlive(pid) {
			time.Sleep(time.Second)
		}
		close(stopped)
	}()
	select {
	case <-stopped:
		log.Printf("Stop the process success.")
	case <-time.After(time.Duration(stopWaitSecs) * time.Second):
		log.Printf("Stop the process timeout, force to kill.")
		syscall.Kill(pid, syscall.SIGKILL)
	}
	return true
}

// Action is the system action.
type Action struct {
	Name string `json:"name"`
}

// ActionStatus is the status of action.
type ActionStatus struct {
	Status bool   `json:"status"`
	Text   string `json:"text"`
}

// Handle plays action.
func (s *System) Handle(action *Action) *ActionStatus {
	var (
		st  = &ActionStatus{}
		err error
	)
	switch action.Name {
	case "status":
		fallthrough
	case "is_alive":
		_, st.Status = s.IsAlive()
	case "start":
		if st.Status, err = s.Start(); err != nil {
			st.Text = err.Error()
		}
	case "stop":
		st.Status = s.Stop()
	case "restart":
		if st.Status, err = s.Restart(); err != nil {
			st.Text = err.Error()
		}
	default:
		st.Status, st.Text = false, fmt.Sprintf("Invalid action: %s", action.Name)
	}
	log.Printf("Action finished")
	return st
}

// WriteFile tries to create parent directory before WriteFile.
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(path.Dir(filename), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, perm)
}

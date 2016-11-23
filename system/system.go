// Copyright 2016 Eleme Inc. All rights reserved.

package system

import (
	"fmt"
	logging "log"
	"os"
)

var (
	log          = logging.New(os.Stdout, "witch: ", logging.Ldate|logging.Ltime|logging.Lmicroseconds|logging.Lshortfile)
	stopWaitSecs = 5
)

// System is the interface of process control system.
type System interface {
	// IsAlive checks process is alive.
	IsAlive() (int, bool)
	// Start starts process.
	Start() (bool, error)
	// Start restart process.
	Restart() (bool, error)
	// Stop stops process.
	Stop() bool
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

// Controller controls the System.
type Controller struct {
	System
}

// Handle plays action.
func (c *Controller) Handle(action *Action) *ActionStatus {
	var (
		st  = &ActionStatus{}
		err error
	)
	switch action.Name {
	case "status":
		fallthrough
	case "is_alive":
		_, st.Status = c.IsAlive()
	case "start":
		if st.Status, err = c.Start(); err != nil {
			st.Text = err.Error()
		}
	case "stop":
		st.Status = c.Stop()
	case "restart":
		if st.Status, err = c.Restart(); err != nil {
			st.Text = err.Error()
		}
	default:
		st.Status, st.Text = false, fmt.Sprintf("Invalid action: %s", action.Name)
	}
	log.Printf("Action finished")
	return st
}

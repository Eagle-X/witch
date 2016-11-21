// Copyright 2016 Eleme Inc. All rights reserved.

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/martini-contrib/render"
)

var (
	// ErrServerError is internal server error.
	ErrServerError = errors.New("Internal Server Error")
	// ErrBadRequest is bad request error.
	ErrBadRequest = errors.New("Bad Request")
)

func sysAction(sys *System, req *http.Request, r render.Render) {
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Read request body error: %s", err)
		r.JSON(http.StatusInternalServerError, ErrServerError)
		return
	}
	log.Printf("Request action: %s", bs)
	action := &Action{}
	if err := json.Unmarshal(bs, action); err != nil {
		log.Printf("Invalid action format: %s", err)
		r.JSON(http.StatusBadRequest, ErrBadRequest)
		return
	}
	r.JSON(http.StatusOK, sys.Handle(action))
}

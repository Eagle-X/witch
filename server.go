// Copyright 2016 Eleme Inc. All rights reserved.

package main

import (
	"net/http"
	"strings"

	"github.com/Eagle-X/witch/system"
	"github.com/braintree/manners"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
	"github.com/martini-contrib/render"
)

// Server is the system RESTful web server.
type Server struct {
	addr string
	m    *martini.ClassicMartini
}

// NewServer inits a system RESTful web server.
func NewServer(addr string, control *system.Controller, cfg *Config) *Server {
	ser := &Server{
		addr: addr,
		m:    martini.Classic(),
	}
	authFunc := auth.BasicFunc(func(username, password string) bool {
		pwd, ok := cfg.Auth[username]
		return ok && pwd == password
	}).(func(http.ResponseWriter, *http.Request, martini.Context))
	ser.m.Map(control)
	ser.m.Use(authInclusive("/api", authFunc))
	ser.m.Use(render.Renderer(render.Options{}))
	ser.m.Put("/api/app/actions", sysAction)
	return ser
}

// Start starts the server.
func (ser *Server) Start() error {
	log.Printf("System webapp start at %s", ser.addr)
	return manners.ListenAndServe(ser.addr, ser.m)
}

// Stop stops the server.
func (ser *Server) Stop() {
	manners.Close()
}

func authInclusive(urlPrefix string, authFunc func(http.ResponseWriter, *http.Request, martini.Context)) martini.Handler {
	return func(resp http.ResponseWriter, req *http.Request, ctx martini.Context) {
		if strings.HasPrefix(req.URL.String(), urlPrefix) {
			if auth := req.URL.Query().Get("auth"); auth != "" && req.Header.Get("Authorization") == "" {
				req.Header.Set("Authorization", "Basic "+auth)
			}
			authFunc(resp, req, ctx)
		} else {
			ctx.Next()
		}
	}
}

package server

import (
	"net/http"
	"time"
)

const (
	// DefaultName is the default name of the server.
	DefaultName = "media-server"
)

var (
	// Version defines the version of server.
	Version = "1.0.0+git"

	// GitHash will be set during make.
	GitHash = "Not provided (use make build instead of go build)"
	// BuildTS will be set during make.
	BuildTS = "Not provided (use make build instead of go build)"

	// health status.
	healthy int32

	// dialClient is a simple HTTP Client.
	dialClient = &http.Client{
		Timeout: time.Second * 10,
	}

	// svr contains server instance.
	svr *Server
)

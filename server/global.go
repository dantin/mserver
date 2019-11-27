package server

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// DefaultName is the default name of the server.
	DefaultName = "media-server"

	// defaultListenAddr is the default addr of the server
	defaultListenAddr = "127.0.0.1:9527"
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

func printVersionInfo() {
	fmt.Printf("%s Version: %s\n", DefaultName, Version)
	fmt.Printf("Git Commit Hash: %s\n", GitHash)
	fmt.Printf("Build Timestamp: %s\n", BuildTS)
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("Go OS/Arch: %s:%s\n", runtime.GOOS, runtime.GOARCH)
}

func showVersionInfo() {
	log.Infof("%s Version: %s", DefaultName, Version)
	log.Infof("Git Commit Hash: %s", GitHash)
	log.Infof("Build Timestamp: %s", BuildTS)
	log.Infof("Go Version: %s", runtime.Version())
	log.Infof("Go OS/Arch: %s:%s", runtime.GOOS, runtime.GOARCH)
}

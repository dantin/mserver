package main

import (
	"log"
	"os"

	"github.com/dantin/mserver/server"
)

func main() {
	cfg := server.NewConfig()
	if err := cfg.Parse(os.Args[1:]); err != nil {
		log.Fatal("bad config file: %v", err)
	}

	svr := server.NewServer(cfg)
	if err := svr.Run(); err != nil {
		log.Fatal("fail to run: %v", err)
	}
}

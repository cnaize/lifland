package main

import (
	"flag"
	"time"

	"github.com/cnaize/lifland/db"
	"github.com/cnaize/lifland/server"
)

var (
	syncDelay time.Duration
)

func init() {
	flag.DurationVar(&syncDelay, "sync-delay", time.Duration(1*time.Second), "sync funds delay")
}

func main() {
	flag.Parse()
	s := server.NewServer(db.NewDB(), syncDelay)
	panic(s.Run("8000"))
}

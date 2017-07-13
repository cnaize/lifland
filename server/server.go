package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cnaize/lifland/db"
	h "github.com/cnaize/lifland/server/handle"
)

type Server struct {
	dbi       db.Interface
	syncDelay time.Duration
	mux       *http.ServeMux
}

func NewServer(dbi db.Interface, syncDelay time.Duration) *Server {
	mux := http.NewServeMux()

	// common
	mux.HandleFunc("/reset", h.Log(h.Reset(dbi)))

	// player
	mux.HandleFunc("/balance", h.Log(h.Balance(dbi)))
	mux.HandleFunc("/take", h.Log(h.Take(dbi)))
	mux.HandleFunc("/fund", h.Log(h.Fund(dbi)))

	// tournament
	mux.HandleFunc("/announceTournament", h.Log(h.Announce(dbi)))
	mux.HandleFunc("/joinTournament", h.Log(h.Join(dbi)))
	mux.HandleFunc("/resultTournament", h.Log(h.Result(dbi)))

	return &Server{
		dbi:       dbi,
		syncDelay: syncDelay,
		mux:       mux,
	}
}

func (s *Server) Run(port string) error {
	fmt.Printf("Server run on port: %s\n", port)
	defer func() {
		fmt.Println("Server stopped")
	}()

	go s.syncFunds()
	return http.ListenAndServe(":"+port, s.mux)
}

func (s *Server) syncFunds() {
	for {
		s.dbi.SyncFunds()
		time.Sleep(s.syncDelay)
	}
}

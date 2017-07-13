package handle

import (
	"net/http"

	"github.com/cnaize/lifland/db"
)

func initTestMux(dbi db.Interface) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/reset", Log(Reset(dbi)))
	mux.HandleFunc("/balance", Log(Balance(dbi)))
	mux.HandleFunc("/take", Log(Take(dbi)))
	mux.HandleFunc("/fund", Log(Fund(dbi)))
	mux.HandleFunc("/announceTournament", Log(Announce(dbi)))
	mux.HandleFunc("/joinTournament", Log(Join(dbi)))
	mux.HandleFunc("/resultTournament", Log(Result(dbi)))
	return mux
}

package handle

import (
	"fmt"
	"net/http"

	"github.com/cnaize/lifland/db"
	"github.com/cnaize/lifland/model"
)

func Announce(dbi db.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := handleAnnounceIn(w, r, dbi)
		if err != nil {
			fmt.Printf("ERROR: Announce(): can't handle input: %+v\n", err)
			return
		}
		tournament := model.NewTournament(in.TournamentId, in.Deposit)
		if err = dbi.AddTournament(tournament); err != nil {
			fmt.Printf("ERROR: Announce(): can't add tournament %d: %+v\n",
				tournament.Id(), err)
			http.Error(w, "", http.StatusConflict)
			return
		}
	}
}

func Join(dbi db.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := handleJoinIn(w, r, dbi)
		if err != nil {
			fmt.Printf("ERROR: Join(): can't handle input: %+v\n", err)
			return
		}
		fund, err := MakeFund(in.Backers, -in.Tournament.Deposit(), dbi)
		if err != nil {
			fmt.Printf("ERROR: Join(): player %s can't make fund for tournament %d: %+v\n",
				in.PlayerId, in.Tournament.Id(), err)
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		}
		if err := in.Tournament.AddPlayer(in.PlayerId, fund); err != nil {
			fmt.Printf("ERROR: Join(): can't add player %s to tournament %d: %+v\n",
				in.PlayerId, in.Tournament.Id(), err)
			http.Error(w, "", http.StatusConflict)
			return
		}
	}
}

func Result(dbi db.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := handleResultIn(w, r, dbi)
		if err != nil {
			fmt.Printf("ERROR: Result(): can't parse input: %+v\n", err)
			return
		}
		funds, err := in.Tournament.Close()
		if err != nil {
			fmt.Printf("ERROR: Result(): can't close tournament %d: %+v\n",
				in.Tournament.Id(), err)
			http.Error(w, "", http.StatusConflict)
			return
		}
		for winner, prize := range in.Winners {
			var players []*model.Player
			for player, _ := range funds[winner.Id()] {
				players = append(players, player)
			}
			MakeFund(players, prize, dbi)
		}
	}
}

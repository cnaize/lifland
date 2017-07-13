package handle

import (
	"fmt"
	"net/http"

	"github.com/cnaize/lifland/db"
	"github.com/cnaize/lifland/model"
	"github.com/cnaize/lifland/util"
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
		fund, err := makeFund(in.Backers, -in.Tournament.Deposit(), dbi)
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
			makeFund(players, prize, dbi)
		}
	}
}

func makeFund(players []*model.Player, points float64, dbi db.Interface) (model.Fund, error) {
	fund := model.Fund{}
	perPlayer := util.Round(points / util.Round(float64(len(players))))
	for i, pl := range players {
		income := perPlayer
		// take rest of points from the player
		// NOTE: the player placed in last position
		if i == len(players)-1 {
			income = points - float64(len(players)-1)*perPlayer
		}

		if err := pl.IncrBalance(income); err != nil {
			fmt.Printf("MakeFund(): can't increase player %s balance by %f points: %+v\n",
				pl.Id(), points, err)
			break
		}
		fund[pl] = income
	}

	// revert if necessary
	if len(fund) != len(players) {
		for pl, income := range fund {
			if err := pl.IncrBalance(-income); err == nil {
				delete(fund, pl)
			}
		}
		// if not all points reverted
		if len(fund) != 0 {
			// store in db for delayed revert
			dbi.AddFund(fund)
		}
		return nil, fmt.Errorf("MakeFund(): failed")
	}
	return fund, nil
}

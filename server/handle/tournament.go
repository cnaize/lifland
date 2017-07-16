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
				tournament.GetId(), err)
			http.Error(w, "", http.StatusConflict)
			return
		}
		dbi.Dump()
	}
}

func Join(dbi db.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := handleJoinIn(w, r, dbi)
		if err != nil {
			fmt.Printf("ERROR: Join(): can't handle input: %+v\n", err)
			return
		}
		fund, err := makeFund(in.Backers, -in.Tournament.GetDeposit(), dbi)
		defer func() {
			if fund == nil {
				return
			}
			dbi.AddFund(fund.Invert())
		}()
		if err != nil {
			fmt.Printf("ERROR: Join(): player %s can't make fund for tournament %d: %+v\n",
				in.PlayerId, in.Tournament.GetId(), err)
			http.Error(w, "", http.StatusUnprocessableEntity)
		} else if err := in.Tournament.AddPlayer(in.PlayerId, fund); err != nil {
			fmt.Printf("ERROR: Join(): can't add player %s to tournament %d: %+v\n",
				in.PlayerId, in.Tournament.GetId(), err)
			http.Error(w, "", http.StatusConflict)
		} else {
			dbi.Dump()
			fund = nil
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
				in.Tournament.GetId(), err)
			http.Error(w, "", http.StatusConflict)
			return
		}
		for winnerId, prize := range in.Winners {
			var playerIds []string
			for playerId, _ := range funds[winnerId] {
				playerIds = append(playerIds, playerId)
			}
			if fund, err := makeFund(playerIds, prize, dbi); err != nil {
				dbi.AddFund(fund.Invert())
			}
		}
		dbi.Dump()
	}
}

func makeFund(playerIds []string, points float64, dbi db.Interface) (model.Fund, error) {
	fund := model.Fund{}
	perPlayer := util.Round(points / util.Round(float64(len(playerIds))))
	for i, id := range playerIds {
		income := perPlayer
		// take rest of points from the player
		// NOTE: the player placed in last position
		if i == len(playerIds)-1 {
			income = points - float64(len(playerIds)-1)*perPlayer
		}

		player := dbi.GetPlayer(id)
		if player == nil {
			fmt.Printf("ERROR: makeFund(): player %s not found", id)
			break
		}
		if err := player.IncrBalance(income); err != nil {
			fmt.Printf("makeFund(): can't increase player %s balance by %f points: %+v\n",
				id, points, err)
			break
		}
		fund[id] = income
	}

	if len(fund) != len(playerIds) {
		return nil, fmt.Errorf("makeFund(): failed")
	}
	return fund, nil
}

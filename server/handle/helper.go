package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cnaize/lifland/db"
	"github.com/cnaize/lifland/model"
)

type inAnnounce struct {
	TournamentId int
	Deposit      float64
}

type inJoin struct {
	Tournament *model.Tournament
	PlayerId   string
	Backers    []string
}

type inResult struct {
	Tournament *model.Tournament
	Winners    model.Fund
}

func handleAnnounceIn(w http.ResponseWriter, r *http.Request, dbi db.Interface) (*inAnnounce, error) {
	query := r.URL.Query()
	tid, err := strconv.Atoi(query.Get("tournamentId"))
	deposit, e := strconv.ParseFloat(query.Get("deposit"), 64)
	if err != nil || e != nil || deposit <= 0 {
		http.Error(w, "", http.StatusBadRequest)
		return nil, fmt.Errorf("invalid input: %v", r.URL.RawQuery)
	}
	return &inAnnounce{
		TournamentId: tid,
		Deposit:      deposit,
	}, nil
}

func handleJoinIn(w http.ResponseWriter, r *http.Request, dbi db.Interface) (*inJoin, error) {
	query := r.URL.Query()
	tid, err := strconv.Atoi(query.Get("tournamentId"))
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return nil, fmt.Errorf("can't parse tournament id %s: %+v", query.Get("tournamentId"), err)
	}
	tournament := dbi.GetTournament(tid)
	if tournament == nil {
		http.Error(w, "", http.StatusNotFound)
		return nil, fmt.Errorf("tournament %d not found", tid)
	}
	// NOTE: the player placed in last position
	qplayers := append(query["backerId"], query.Get("playerId"))
	for i, backerId := range qplayers {
		// check duplicates
		for j, id := range qplayers {
			if i != j && id == backerId {
				// duplicate found
				http.Error(w, "", http.StatusBadRequest)
				return nil, fmt.Errorf("passed duplicated player %s to tournament %d", backerId, tid)
			}
		}

		backer := dbi.GetPlayer(backerId)
		if backer == nil {
			http.Error(w, "", http.StatusNotFound)
			return nil, fmt.Errorf("player %s not found", backerId)
		}
	}
	return &inJoin{
		Tournament: tournament,
		PlayerId:   query.Get("playerId"),
		Backers:    qplayers,
	}, nil
}

func handleResultIn(w http.ResponseWriter, r *http.Request, dbi db.Interface) (*inResult, error) {
	type inData struct {
		Winners []struct {
			PlayerId string  `json:"playerId"`
			Prize    float64 `json:"prize"`
		} `json:"winners"`
	}

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return nil, fmt.Errorf("invalid method %s", r.Method)
	}
	var in inData
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return nil, fmt.Errorf("can't parse in json: %+v", err)
	}
	tournament := dbi.GetOldestTournament()
	if tournament == nil {
		http.Error(w, "", http.StatusNotFound)
		return nil, fmt.Errorf("tournament not found")
	}
	winners := make(model.Fund)
	for i, winner := range in.Winners {
		for j, wnr := range in.Winners {
			if i != j && winner == wnr {
				http.Error(w, "", http.StatusBadRequest)
				return nil, fmt.Errorf("passed duplicated player %s to tournament %d",
					winner.PlayerId, tournament.GetId())
			}
		}

		if !tournament.HasPlayer(winner.PlayerId) {
			http.Error(w, "", http.StatusBadRequest)
			return nil, fmt.Errorf("player %s not in tournament %d",
				winner.PlayerId, tournament.GetId())
		}
		if winner.Prize <= 0 {
			http.Error(w, "", http.StatusBadRequest)
			return nil, fmt.Errorf("invalid prize %f for player %s", winner.Prize, winner.PlayerId)
		}
		player := dbi.GetPlayer(winner.PlayerId)
		if player == nil {
			http.Error(w, "", http.StatusNotFound)
			return nil, fmt.Errorf("player %s not found", winner.PlayerId)
		}
		winners[winner.PlayerId] = winner.Prize
	}
	return &inResult{
		Tournament: tournament,
		Winners:    winners,
	}, nil
}

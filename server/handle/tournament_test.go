package handle

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/cnaize/lifland/db"
	"github.com/cnaize/lifland/model"
	"github.com/cnaize/lifland/util"
)

func TestTournamentAnnounce(t *testing.T) {
	makeUri := func(id, deposit string) string {
		return fmt.Sprintf("/announceTournament?tournamentId=%s&deposit=%s", id, deposit)
	}

	tests := []struct {
		tournId     string
		deposit     string
		wantCode    int
		wantDeposit float64
	}{
		{"", "", http.StatusBadRequest, 0},
		{"1", "", http.StatusBadRequest, 0},
		{"", "10", http.StatusBadRequest, 0},
		{"qwe", "10", http.StatusBadRequest, 0},
		{"1", "qwe", http.StatusBadRequest, 0},
		{"1", "10", http.StatusConflict, 0},
		{"2", "20", http.StatusOK, 20},
	}

	for _, test := range tests {
		tourn1 := model.NewTournament(1, 10)
		dbi := db.NewDB()
		dbi.AddTournament(tourn1)

		uri := makeUri(test.tournId, test.deposit)
		r, _ := http.NewRequest(http.MethodPost, uri, nil)
		w := httptest.NewRecorder()
		initTestMux(dbi).ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("invalid code %d for uri %s", w.Code, uri)
		}
		id, _ := strconv.Atoi(test.tournId)
		tournament := dbi.GetTournament(id)
		if test.wantCode != http.StatusOK || tournament == nil {
			continue
		}
		if util.Round(tournament.Deposit()) != util.Round(test.wantDeposit) {
			t.Errorf("invalid deposit for uri %s: want %f, got %f",
				uri, test.wantDeposit, tournament.Deposit())
		}
	}
}

func TestTournamentJoin(t *testing.T) {
	makeUri := func(id, playerId string, backers ...string) string {
		uri := fmt.Sprintf("/joinTournament?tournamentId=%s&playerId=%s",
			id, playerId)
		for _, backer := range backers {
			uri = fmt.Sprintf("%s&backerId=%s", uri, backer)
		}
		return uri
	}

	tests := []struct {
		tournId  string
		playerId string
		backers  []string
		wantCode int
	}{
		{"", "", []string{}, http.StatusBadRequest},
		{"1", "invalid", []string{}, http.StatusNotFound},
		{"qwe", "10", []string{}, http.StatusBadRequest},
		{"1", "10", []string{"10"}, http.StatusBadRequest},
		{"1", "10", []string{"invalid"}, http.StatusNotFound},
		{"1", "10", []string{}, http.StatusOK},
		{"1", "10", []string{"20"}, http.StatusOK},
		{"2", "10", []string{"20"}, http.StatusUnprocessableEntity},
	}

	for _, test := range tests {
		player10 := model.NewPlayer("10")
		player20 := model.NewPlayer("20")
		player10.IncrBalance(10)
		player20.IncrBalance(20)
		tourn1 := model.NewTournament(1, 10)
		tourn2 := model.NewTournament(2, 30)
		dbi := db.NewDB()
		dbi.AddPlayer(player10)
		dbi.AddPlayer(player20)
		dbi.AddTournament(tourn1)
		dbi.AddTournament(tourn2)

		uri := makeUri(test.tournId, test.playerId, test.backers...)
		r, _ := http.NewRequest(http.MethodPost, uri, nil)
		w := httptest.NewRecorder()
		initTestMux(dbi).ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("invalid code %d for uri %s", w.Code, uri)
		}
		id, _ := strconv.Atoi(test.tournId)
		tournament := dbi.GetTournament(id)
		if test.wantCode != http.StatusOK || tournament == nil {
			continue
		}
		if test.playerId != player10.Id() ||
			(len(test.backers) > 0 && test.backers[0] != player20.Id()) {
			continue
		}
		if len(test.backers) == 0 {
			if player10.Balance() != 0.0 {
				t.Errorf("invalid balance for uri %s: want %f, got %f",
					uri, 0.0, player10.Balance())
			}
		} else if len(test.backers) == 1 {
			if player10.Balance() != 5.0 {
				t.Errorf("invalid balance for uri %s: want %f, got %f",
					uri, 5.0, player10.Balance())
			}
			if player20.Balance() != 15.0 {
				t.Errorf("invalid balance for uri %s: want %f, got %f",
					uri, 15.0, player20.Balance())
			}
		}
	}
}

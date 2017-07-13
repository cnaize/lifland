package handle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/cnaize/lifland/db"
	"github.com/cnaize/lifland/model"
	"github.com/cnaize/lifland/util"
)

func TestPlayerFund(t *testing.T) {
	makeUri := func(id, points string) string {
		return fmt.Sprintf("/fund?playerId=%s&points=%s", id, points)
	}

	tests := []struct {
		playerId    string
		points      string
		wantCode    int
		wantBalance float64
	}{
		{"", "", http.StatusBadRequest, 10},
		{"10", "", http.StatusBadRequest, 10},
		{"", "10", http.StatusBadRequest, 10},
		{"qwe", "10", http.StatusOK, 10},
		{"10", "qwe", http.StatusBadRequest, 10},
		{"10", "0", http.StatusBadRequest, 10},
		{"10", "-10", http.StatusBadRequest, 10},
		{"10", "20", http.StatusOK, 30},
		{"30", "30", http.StatusOK, 30},
	}

	for _, test := range tests {
		player10 := model.NewPlayer("10")
		player10.IncrBalance(10)
		dbi := db.NewDB()
		dbi.AddPlayer(player10)

		uri := makeUri(test.playerId, test.points)
		r, _ := http.NewRequest(http.MethodPost, uri, nil)
		w := httptest.NewRecorder()
		initTestMux(dbi).ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("invalid code %d for uri %s", w.Code, uri)
		}
		player := dbi.GetPlayer(test.playerId)
		if test.wantCode != http.StatusOK || player == nil {
			continue
		}
		if util.Round(player.Balance()) != util.Round(test.wantBalance) {
			t.Errorf("invalid balance for uri %s: want %f, got %f",
				uri, test.wantBalance, player.Balance())
		}
	}
}

func TestPlayerTake(t *testing.T) {
	makeUri := func(id, points string) string {
		return fmt.Sprintf("/take?playerId=%s&points=%s", id, points)
	}

	tests := []struct {
		playerId    string
		points      string
		wantCode    int
		wantBalance float64
	}{
		{"", "", http.StatusNotFound, 0},
		{"10", "", http.StatusBadRequest, 0},
		{"", "10", http.StatusNotFound, 0},
		{"qwe", "10", http.StatusNotFound, 0},
		{"10", "qwe", http.StatusBadRequest, 0},
		{"10", "-10", http.StatusBadRequest, 0},
		{"10", "10.1", http.StatusUnprocessableEntity, 0},
		{"10", "9.9", http.StatusOK, 0.1},
	}

	for _, test := range tests {
		player10 := model.NewPlayer("10")
		player10.IncrBalance(10)
		dbi := db.NewDB()
		dbi.AddPlayer(player10)

		uri := makeUri(test.playerId, test.points)
		r, _ := http.NewRequest(http.MethodPost, uri, nil)
		w := httptest.NewRecorder()
		initTestMux(dbi).ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("invalid code %d for uri %s", w.Code, uri)
		}
		player := dbi.GetPlayer(test.playerId)
		if test.wantCode != http.StatusOK || player == nil {
			continue
		}
		if util.Round(player.Balance()) != util.Round(test.wantBalance) {
			t.Errorf("invalid balance for uri %s: want %f, got %f",
				uri, test.wantBalance, player.Balance())
		}
	}
}

func TestPlayerBalance(t *testing.T) {
	makeUri := func(id string) string {
		return fmt.Sprintf("/balance?playerId=%s", id)
	}

	player10 := model.NewPlayer("10")
	player10.IncrBalance(10)
	player10Data := map[string]interface{}{
		"playerId": "10",
		"balance":  10.0,
	}
	dbi := db.NewDB()
	dbi.AddPlayer(player10)

	tests := []struct {
		playerId string
		wantCode int
		wantData map[string]interface{}
	}{
		{"", http.StatusNotFound, nil},
		{"20", http.StatusNotFound, nil},
		{"10", http.StatusOK, player10Data},
	}

	for _, test := range tests {
		uri := makeUri(test.playerId)
		r, _ := http.NewRequest(http.MethodPost, uri, nil)
		w := httptest.NewRecorder()
		initTestMux(dbi).ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("invalid code %d for uri %s", w.Code, uri)
		}
		if test.wantData != nil {
			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("can't read body: %+v", err)
			}
			var data map[string]interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				t.Errorf("can't unmarshal body: %+v", err)
			}
			if !reflect.DeepEqual(data, test.wantData) {
				t.Errorf("invalid data: want\n%#v\ngot\n%#v", data, test.wantData)
			}
		}
	}
}

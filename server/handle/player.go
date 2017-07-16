package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cnaize/lifland/db"
	"github.com/cnaize/lifland/model"
)

func Take(dbi db.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		player := dbi.GetPlayer(query.Get("playerId"))
		if player == nil {
			fmt.Printf("ERROR: Take(): player %s not found\n", query.Get("playerId"))
			http.Error(w, "", http.StatusNotFound)
			return
		}

		points, err := strconv.ParseFloat(query.Get("points"), 64)
		if err != nil || points <= 0 {
			fmt.Printf("ERROR: Take(): invalid points: %s\n", query.Get("points"))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if err := player.IncrBalance(-points); err != nil {
			fmt.Printf("ERROR: Take(): can't take %f points from player %s: %+v\n",
				points, player.GetId(), err)
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		}
		dbi.Dump()
	}
}

func Fund(dbi db.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		pid := query.Get("playerId")
		points, err := strconv.ParseFloat(query.Get("points"), 64)
		if pid == "" || err != nil || points <= 0 {
			fmt.Printf("ERROR: Fund(): invalid input: %v\n", r.URL.RawQuery)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		player := dbi.GetPlayer(pid)
		if player == nil {
			player = model.NewPlayer(pid)
			if err := dbi.AddPlayer(player); err != nil {
				fmt.Printf("ERROR: Fund(): can't add player: %+v\n", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}
		if err := player.IncrBalance(points); err != nil {
			fmt.Printf("ERROR: Fund(): can't give %f points to player %s: %+v\n",
				points, player.GetId(), err)
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		}
		dbi.Dump()
	}
}

func Balance(dbi db.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		query := r.URL.Query()
		player := dbi.GetPlayer(query.Get("playerId"))
		if player == nil {
			fmt.Printf("ERROR: Balance(): player %s not found\n", query.Get("playerId"))
			http.Error(w, "", http.StatusNotFound)
			return
		}

		data := map[string]interface{}{
			"playerId": player.GetId(),
			"balance":  player.GetBalance(),
		}
		resp, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("ERROR: Balance(): can't marshal data %v: %+v\n", data, err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(resp); err != nil {
			fmt.Printf("ERROR: Balance(): can't write response %s: %+v\n", resp, err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

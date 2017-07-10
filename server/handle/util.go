package handle

import (
	"fmt"
	"net/http"

	"github.com/cnaize/lifland/db"
	"github.com/cnaize/lifland/model"
)

func Log(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("New request: %s %s\n", r.Method, r.RequestURI)
		fn(w, r)
	}
}

func MakeFund(players []*model.Player, points float64, dbi db.Interface) (model.Fund, error) {
	fund := model.Fund{}
	perPlayer := Round(points / float64(len(players)))
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

func Round(f float64) float64 {
	return float64(int64(f*100)) / 100
}

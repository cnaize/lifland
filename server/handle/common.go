package handle

import (
	"net/http"

	"github.com/cnaize/lifland/db"
)

func Reset(dbi db.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbi.Reset()
	}
}

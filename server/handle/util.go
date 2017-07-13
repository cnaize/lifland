package handle

import (
	"fmt"
	"net/http"
)

func Log(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("New request: %s %s\n", r.Method, r.RequestURI)
		fn(w, r)
	}
}

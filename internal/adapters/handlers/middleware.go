package handlers

import (
	"net/http"

	"github.com/petrostrak/xm-companies/utils"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !utils.IsAuthenticated(r) {
			// TODO: error message
			return
		}
		next.ServeHTTP(w, r)
	})
}

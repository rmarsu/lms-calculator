package transport

import (
	"lms-1/internal/domain"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func CheckPOSTMethod(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {

			w.Header().Set("Content-Type", "application/json")

			http.Error(w, domain.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

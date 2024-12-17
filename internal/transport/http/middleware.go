package transport

import (
	"lms-1/internal/domain"
	rjson "lms-1/pkg/json"
	"net/http"
)

func CheckPOSTMethod(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			rjson.SendJson(w, domain.ErrResponse{Error: "Only POST method is allowed"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

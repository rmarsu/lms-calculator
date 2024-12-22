package transport

import (
	"encoding/json"
	"lms-1/internal/domain"
	rjson "lms-1/pkg/json"
	"lms-1/pkg/logger"
	"net/http"
)

func (h *Handler) calculateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var input domain.Input

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Info(err)

		w.WriteHeader(http.StatusBadRequest)
		rjson.SendJson(w, domain.ErrResponse{Error: domain.ErrInvalidJSON})
		return
	}

	result, err := h.s.Calc.Calculate(input.Expression)
	if err != nil {
		logger.Info(err)

		w.WriteHeader(http.StatusUnprocessableEntity)
		rjson.SendJson(w, domain.ErrResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	rjson.SendJson(w, domain.OkResponse{Result: result})
}

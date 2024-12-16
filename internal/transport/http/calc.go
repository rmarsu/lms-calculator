package transport

import (
	"encoding/json"
	"lms-1/internal/domain"
	"lms-1/pkg/logger"
	"net/http"
)

func (h *Handler) calculateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input domain.Input

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Info(err)

		http.Error(w, domain.ErrBadRequest, http.StatusBadRequest)
		return
	}

	result, err := h.s.Calc.Calculate(input.Expression)
	if err != nil {
		logger.Info(err)

		http.Error(w, domain.ErrBadExpression, http.StatusUnprocessableEntity)
		return
	}

	response := domain.OKAnswer{Result: result}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

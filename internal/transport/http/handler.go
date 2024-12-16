package transport

import (
	"lms-1/internal/service"
	"net/http"
)

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/api/v1/calculate", CheckPOSTMethod(h.calculateHandler))

	return router
}

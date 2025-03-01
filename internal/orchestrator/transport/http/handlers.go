package http

import (
	"lms-1/internal/orchestrator/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	api := r.Group("api")
	{
		v1 := api.Group("v1")
		{
			v1.GET("/expressions", h.getExpressions)
			v1.GET("/expressions/:id", h.getExpressionById)

			v1.POST("/calculate", h.addToQueue)
		}
	}
	internal := r.Group("internal")
	{
		internal.GET("/task", h.getTask)
	}

	return r
}

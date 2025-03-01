package http

import (
	"lms-1/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getTask(c *gin.Context) {
	task, ok := h.s.Queue.GetTask()
	if !ok {
		newErrorResponse(http.StatusNotFound, domain.ErrNoTasksAvailable, c)
		return
	}
	c.JSON(http.StatusOK, task)
}

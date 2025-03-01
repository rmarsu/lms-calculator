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

func (h *Handler) rollbackResult(c *gin.Context) {
	var input *domain.AgentInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		newErrorResponse(http.StatusBadRequest, domain.ErrInvalidJSON, c)
		return
	}
	err := h.s.Queue.RollbackResult(input.Id, input.Result)
	if err != nil {
		newErrorResponse(http.StatusInternalServerError, err.Error(), c)
		return
	}
	c.Status(http.StatusOK)
}

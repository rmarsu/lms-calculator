package http

import (
	"lms-1/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) getExpressions(c *gin.Context) {
	exprs := h.s.Queue.GetExpressions()
	c.JSON(http.StatusOK, exprs)
}

func (h *Handler) getExpressionById(c *gin.Context) {
	id := c.Param("id")
	expr, ok := h.s.Queue.GetExpressionById(id)
	if !ok {
		newErrorResponse(http.StatusNotFound, domain.ErrExpressionNotFound, c)
		return
	}

	c.JSON(http.StatusOK, expr)
}

func (h *Handler) addToQueue(c *gin.Context) {
	var input *domain.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(http.StatusBadRequest, domain.ErrInvalidJSON, c)
		return
	}
	expr := &domain.Expression{
		Id:         uuid.New().String(),
		Expression: input.Expression,
		Status:     domain.StatusPending,
		Result:     0.0,
	}
	h.s.Queue.AddExpression(*expr)
	go func() {
		if err := h.s.Queue.RunTask(expr.Id); err != nil {
			h.s.Queue.RemoveExpression(expr.Id)
			newErrorResponse(http.StatusUnprocessableEntity, err.Error(), c)
			return
		}
	}()

	c.JSON(http.StatusCreated, expr)
}

package queue

import "lms-1/internal/domain"

type Timings struct {
	TimeAdditionMs       int32
	TimeSubtractionMs    int32
	TimeMultiplicationMs int32
	TimeDivisionMs       int32
}

type Queue struct {
	Expressions map[string]*domain.Expression
	Timings
}

func New(t Timings) *Queue {
	return &Queue{
		Expressions: make(map[string]*domain.Expression),
		Timings:     t,
	}
}

func (q *Queue) AddExpression(expr domain.Expression) {
	q.Expressions[expr.Id] = &expr
}

func (q *Queue) GetExpressionById(id string) (*domain.Expression, bool) {
	expr, ok := q.Expressions[id]
	return expr, ok
}

func (q *Queue) GetExpressions() []domain.Expression {
	var exprs []domain.Expression
	for _, expr := range q.Expressions {
		exprs = append(exprs, *expr)
	}
	return exprs
}

func (q *Queue) GetExpression() *domain.Expression {
	for _, expr := range q.Expressions {
		return expr
	}
	return &domain.Expression{}
}

func (q *Queue) GetTask() (*domain.Task, bool) {
	expr := q.GetExpression()
	expr.Status = domain.StatusInProgress

	return &domain.Task{
		Id:              expr.Id,
		Arg1:            0.0, // TODO
		Arg2:            0.0, // TODO
		Operation:       "+", // TODO
		OperationTimeMs: q.TimeAdditionMs,
	}, true
}

func (q *Queue) RemoveExpression(id string) {
	delete(q.Expressions, id)
}

func (q *Queue) Size() int {
	return len(q.Expressions)
}

package domain

const (
	ErrBadRequest          = `{"error":"Bad request"}`
	ErrBadExpression       = `{"error":"Expression is not valid"}`
	ErrInternalServerError = `{"error":"Internal server error"}`
	ErrMethodNotAllowed    = `{"error":"Only POST method is allowed"}`
)

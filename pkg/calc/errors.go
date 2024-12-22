package calc

import "errors"

var (
	ErrInvalidParentheses = errors.New("неверно расставленные скобки")
	ErrInvalidExpression = errors.New("некорректное выражение")
	ErrDivisionByZero    = errors.New("деление на 0")
)

package domain

type ErrResponse struct {
	Error string `json:"error"`
}

const (
	ErrInvalidMethod      = "Доступен только метод POST"
	ErrInvalidExpression  = "невалидное выражение"
	ErrInvalidJSON        = "Ошибка при парсинге JSON"
	ErrExpressionNotFound = "Извините, выражение не найдено"
	ErrNoTasksAvailable   = "Нет доступных задач"
	ErrTaskNotFound       = "задача не найдена"
)

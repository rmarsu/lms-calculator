package domain

type ErrResponse struct {
	Error string `json:"error"`
}

const (
	ErrInvalidMethod      = "Доступен только метод POST"
	ErrInvalidJSON        = "Ошибка при парсинге JSON"
	ErrExpressionNotFound = "Извините, выражение не найдено"
	ErrNoTasksAvailable   = "Нет доступных задач"
)

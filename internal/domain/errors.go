package domain

type ErrResponse struct {
	Error string `json:"error"`
}

const (
	ErrInvalidMethod = "Доступен только метод POST"
	ErrInvalidJSON = "Ошибка при парсинге JSON"
)
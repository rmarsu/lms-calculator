package domain

type OKAnswer struct {
	Result float64 `json:"result"`
}

type ErrorAnswer struct {
	Error string `json:"error"`
}

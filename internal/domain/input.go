package domain

type UserInput struct {
	Expression string `json:"expression"`
}

type AgentInput struct {
	Id     string  `json:"id"`
	Result float64 `json:"result"`
}

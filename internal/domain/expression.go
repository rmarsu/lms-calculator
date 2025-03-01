package domain

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
)

type Expression struct {
	Id         string `json:"id"`
	Expression string
	Status     Status  `json:"status"`
	Result     float64 `json:"result"`
}

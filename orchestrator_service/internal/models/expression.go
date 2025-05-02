package models

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
)

type Expression struct {
	Id         int64   `json:"id"`
	Expression string  `json:"-"`
	Status     Status  `json:"status"`
	Result     float64 `json:"result"`
}

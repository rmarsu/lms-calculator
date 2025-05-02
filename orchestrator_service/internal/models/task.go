package models

import "time"

type Task struct {
	Id            int64        `db:"id"`
	Arg1          float64      `db:"arg1"`
	Arg2          float64      `db:"arg2"`
	Operation     string       `db:"operation"`
	OperationTime int64        `db:"operation_time_ms"`
	ExprID        int64        `db:"expr_id"`
	CreatedAt     time.Time    `db:"created_at"`
}

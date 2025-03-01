package domain

type Task struct {
	Id              string       `json:"id"`
	Arg1            float64      `json:"arg1"`
	Arg2            float64      `json:"arg2"`
	Operation       string       `json:"operation"`
	OperationTimeMs int32        `json:"operation_time"`
	ResultChan      chan float64 `json:"-"`
}

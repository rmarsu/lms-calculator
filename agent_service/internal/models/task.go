package models

type Task struct {
	Id       int64
	Arg1     float32
	Arg2     float32
	Op       string
	OpTimeMs int64
}

package service

import "lms-1/internal/orchestrator/queue"

type Service struct {
	Queue *queue.Queue
	Deps  *Deps
}

type Deps struct {
	TimeAdditionMs       int32
	TimeSubtractionMs    int32
	TimeMultiplicationMs int32
	TimeDivisionMs       int32
}

func New(deps *Deps) *Service {
	return &Service{
		Queue: queue.New(queue.Timings{
			TimeAdditionMs:       deps.TimeAdditionMs,
			TimeSubtractionMs:    deps.TimeSubtractionMs,
			TimeMultiplicationMs: deps.TimeMultiplicationMs,
			TimeDivisionMs:       deps.TimeDivisionMs,
		}),
	}
}

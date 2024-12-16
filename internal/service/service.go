package service

import "lms-1/pkg/calc"

type Service struct {
	Calc *CalcService
}

type Deps struct {
	Calc *calc.Calc
}

func NewService(deps *Deps) *Service {
	return &Service{
		Calc: NewCalcService(deps.Calc),
	}
}

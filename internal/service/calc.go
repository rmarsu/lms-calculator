package service

import "lms-1/pkg/calc"

type CalcService struct {
	calc *calc.Calc
}

func NewCalcService(calc *calc.Calc) *CalcService {
	return &CalcService{
		calc: calc,
	}
}

func (s *CalcService) Calculate(expr string) (float64, error) {
	return s.calc.Calculate(expr)
}

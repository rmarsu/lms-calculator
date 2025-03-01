package calc

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func evaluateAST(node ast.Node) (float64, error) {
	switch n := node.(type) {
	case *ast.BasicLit:
		val, _ := strconv.ParseFloat(n.Value, 64)
		return val, nil
	case *ast.BinaryExpr:
		left, err := evaluateAST(n.X)
		if err != nil {
			return 0, err
		}
		right, err := evaluateAST(n.Y)
		if err != nil {
			return 0, err
		}
		switch n.Op {
		case token.ADD:
			return left + right, nil
		case token.SUB:
			return left - right, nil
		case token.MUL:
			return left * right, nil
		case token.QUO:
			if right == 0 {
				return 0, errors.New("division by zero")
			}
			return left / right, nil
		}
	}
	return 0, errors.New("invalid expression")
}

func ParseAndEvaluate(expression string) (float64, error) {
	node, err := parser.ParseExpr(expression)
	if err != nil {
		return 0, err
	}
	return evaluateAST(node)
}

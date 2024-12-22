package calc

import (
	"strconv"
	"strings"
	"unicode"
)

type Calc struct{}

func NewCalc() *Calc {
	return &Calc{}
}

func (c *Calc) Calculate(expression string) (float64, error) {
	tokens := tokenize(expression)
	rpn, err := shuntingYard(tokens)
	if err != nil {
		return 0, err
	}
	return evaluateRPN(rpn)
}

func tokenize(expression string) []string {
	var tokens []string
	var number strings.Builder

	for i, char := range expression {
		if unicode.IsDigit(char) || char == '.' {
			number.WriteRune(char)
		} else {
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}
			if !unicode.IsSpace(char) {
				if char == '-' && (i == 0 || isOperator(rune(expression[i-1])) || expression[i-1] == '(') {
					number.WriteRune(char)
				} else {
					tokens = append(tokens, string(char))
				}
			}
		}
	}
	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}
	return tokens
}

func isOperator(char rune) bool {
	return char == '+' || char == '-' || char == '*' || char == '/'
}

func shuntingYard(tokens []string) ([]string, error) {
	var output []string
	var stack []string

	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	for _, token := range tokens {
		switch {
		case isNumber(token):
			output = append(output, token)
		case token == "(":
			stack = append(stack, token)
		case token == ")":
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, ErrInvalidParentheses
			}
			stack = stack[:len(stack)-1]
		default:
			for len(stack) > 0 && precedence[stack[len(stack)-1]] >= precedence[token] {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, ErrInvalidParentheses
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func evaluateRPN(rpn []string) (float64, error) {
	var stack []float64

	for _, token := range rpn {
		if isNumber(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, ErrInvalidExpression
			}
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			var result float64
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, ErrDivisionByZero
				}
				result = a / b
			}
			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpression
	}

	return stack[0], nil
}

func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

package calc

import (
	"errors"
	"go/ast"
	"go/token"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		node     ast.Node
		expected error
	}{
		{&ast.BasicLit{Kind: token.INT, Value: "1"}, nil},
		{&ast.BinaryExpr{
			X:  &ast.BasicLit{Kind: token.INT, Value: "1"},
			Y:  &ast.BasicLit{Kind: token.INT, Value: "2"},
			Op: token.ADD,
		}, nil},
		{&ast.BinaryExpr{
			X:  &ast.BasicLit{Kind: token.INT, Value: "1"},
			Y:  &ast.BasicLit{Kind: token.INT, Value: "2"},
			Op: token.SUB,
		}, nil},
		{&ast.ParenExpr{
			X: &ast.BinaryExpr{
				X:  &ast.BasicLit{Kind: token.INT, Value: "1"},
				Y:  &ast.BasicLit{Kind: token.INT, Value: "2"},
				Op: token.MUL,
			},
		}, nil},

		{&ast.BinaryExpr{
			X:  &ast.BasicLit{Kind: token.INT, Value: "1"},
			Y:  &ast.BasicLit{Kind: token.INT, Value: "2"},
			Op: token.STRUCT, 
		}, errors.New("invalid binary operator")},
		{&ast.CallExpr{}, errors.New("invalid expression")},
	}

	for _, test := range tests {
		err := validate(test.node)
		if (err != nil) != (test.expected != nil) {
			t.Errorf("validate(%v) = %v; want %v", test.node, err, test.expected)
			continue
		}
		if err != nil && err.Error() != test.expected.Error() {
			t.Errorf("validate(%v) = %v; want %v", test.node, err, test.expected)
		}
	}
}

func TestParseAst(t *testing.T) {
	tests := []struct {
		expression string
		expected   bool
	}{
		{"1 + 2", true},
		{"3 - 4", true},
		{"5 * (6 / 7)", true},
		{"1 + a", false},
		{"1 ^ 2", false},
	}

	for _, test := range tests {
		node, err := ParseAst(test.expression)

		if test.expected && err != nil {
			t.Errorf("ParseAst(%q) failed; got error %v", test.expression, err)
		}

		if !test.expected && err == nil {
			t.Errorf("ParseAst(%q) succeeded; expected an error", test.expression)
		}

		if test.expected && node == nil {
			t.Errorf("ParseAst(%q) returned nil node; expected non-nil", test.expression)
		}
	}
}

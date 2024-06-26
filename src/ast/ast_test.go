package ast

import (
	"testing"

	"github.com/tjapit/monkey/src/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "anotherVar",
					},
					Value: "anotherVar",
				},
			},
		},
	}

	expected := "let myVar = anotherVar;"
	if program.String() != expected {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

// debugging nil-pointer for (*CallExpression).String()
// turns out, *bytes.Buffer was the problem. Changed to bytes.Buffer instead of
// pointer fixes the problem.
func TestCallExpressionStatement(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "add",
				},
				Expression: &CallExpression{
					Token: token.Token{
						Type:    token.LPAREN,
						Literal: "(",
					},
					Function: &Identifier{
						Token: token.Token{
							Type:    token.IDENT,
							Literal: "add",
						},
						Value: "add",
					},
					Arguments: []Expression{},
				},
			},
		},
	}

	if program.String() != "add()" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

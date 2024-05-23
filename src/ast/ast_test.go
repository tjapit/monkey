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

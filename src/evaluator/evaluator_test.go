package evaluator

import (
	"testing"

	"github.com/tjapit/monkey/src/lexer"
	"github.com/tjapit/monkey/src/object"
	"github.com/tjapit/monkey/src/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected int64
	}{
		{
			desc:     "Test 1",
			input:    "5",
			expected: 5,
		},
		{
			desc:     "Test 2",
			input:    "10",
			expected: 10,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			evaluated := testEval(tC.input)
			testIntegerObject(t, evaluated, tC.expected)
		})
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf(
			"object has incorrect value. want=%d, got=%d",
			expected,
			result.Value,
		)
		return false
	}

	return true
}

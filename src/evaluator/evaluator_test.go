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
		{"Test 1", "5", 5},
		{"Test 2", "10", 10},
		{"Test 3", "-5", -5},
		{"Test 4", "-10", -10},
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

func TestEvalBooleanLiteral(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected bool
	}{
		{
			desc:     "Test true",
			input:    "true;",
			expected: true,
		},
		{
			desc:     "Test false",
			input:    "false;",
			expected: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			evaluated := testEval(tC.input)
			testBooleanObject(t, evaluated, tC.expected)
		})
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%t (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf(
			"object has wrong value. want=%t, got=%t",
			expected,
			result.Value,
		)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected bool
	}{
		{"Test 1", "!true", false},
		{"Test 2", "!false", true},
		{"Test 3", "!5", false},
		{"Test 4", "!!true", true},
		{"Test 5", "!!false", false},
		{"Test 6", "!!5", true},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			evaluated := testEval(tC.input)
			testBooleanObject(t, evaluated, tC.expected)
		})
	}
}

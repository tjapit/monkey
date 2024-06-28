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
		{"Test 5", "5 + 5 + 5 + 5 - 10", 10},
		{"Test 6", "2 * 2 * 2 * 2 * 2", 32},
		{"Test 7", "-50 + 100 + -50", 0},
		{"Test 8", "5 * 2 + 10", 20},
		{"Test 9", "5 + 2 * 10", 25},
		{"Test 10", "20 + 2 * -10", 0},
		{"Test 11", "50 / 2 * 2 + 10", 60},
		{"Test 12", "2 * (5 + 10)", 30},
		{"Test 13", "3 * 3 * 3 + 10", 37},
		{"Test 14", "3 * (3 * 3) + 10", 37},
		{"Test 15", "(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
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

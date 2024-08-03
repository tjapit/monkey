package vm

import (
	"fmt"
	"testing"

	"github.com/tjapit/monkey/src/ast"
	"github.com/tjapit/monkey/src/compiler"
	"github.com/tjapit/monkey/src/lexer"
	"github.com/tjapit/monkey/src/object"
	"github.com/tjapit/monkey/src/parser"
)

type vmTestCase struct {
	desc     string
	input    string
	expected interface{}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf(
			"object is not Integer. got=%T (%+v)",
			result,
			result,
		)
	}

	if result.Value != expected {
		return fmt.Errorf(
			"object has wrong value. want=%d, got=%d",
			expected,
			result.Value,
		)
	}

	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf(
			"object has wrong value. want=%t, got=%t",
			expected,
			result.Value,
		)
	}

	return nil
}

func testExpectedObject(
	t *testing.T,
	expected interface{},
	actual object.Object,
) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}

	case bool:
		err := testBooleanObject(bool(expected), actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	}
}

func runVmTests(t *testing.T, testCases []vmTestCase) {
	t.Helper()

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			program := parse(tC.input)

			comp := compiler.New()
			err := comp.Compile(program)
			if err != nil {
				t.Fatalf("compiler error: %s", err)
			}

			vm := New(comp.Bytecode())
			err = vm.Run()
			if err != nil {
				t.Fatalf("vm error: %s", err)
			}

			stackElem := vm.LastPopped()

			testExpectedObject(t, tC.expected, stackElem)
		})
	}
}

func TestIntegerArithmetic(t *testing.T) {
	testCases := []vmTestCase{
		{"Test 1", "1", 1},
		{"Test 2", "2", 2},
		{"Test 3", "1 + 2", 3},
		{"Test 4", "1 - 2", -1},
		{"Test 5", "2 * 2", 4},
		{"Test 6", "2 / 1", 2},
		{"Test 7", "50 / 2 * 2 + 10 - 5", 55},
		{"Test 8", "5 * (2 + 10)", 60},
		{"Test 9", "-5", -5},
		{"Test 10", "-50 + 100 + -50", 0},
		{"Test 11", "(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	runVmTests(t, testCases)
}

func TestBooleanExpressions(t *testing.T) {
	testCases := []vmTestCase{
		{"Test True", "true", true},
		{"Test False", "false", false},
		{"Test 1", "1 < 2", true},
		{"Test 2", "1 > 2", false},
		{"Test 3", "1 > 1", false},
		{"Test 4", "1 < 1", false},
		{"Test 5", "1 == 1", true},
		{"Test 6", "1 != 1", false},
		{"Test 7", "1 == 2", false},
		{"Test 8", "1 != 2", true},
		{"Test 9", "true == true", true},
		{"Test 10", "false == false", true},
		{"Test 11", "true == false", false},
		{"Test 12", "true != false", true},
		{"Test 13", "false != true", true},
		{"Test 14", "(1 < 2) == true", true},
		{"Test 15", "(1 < 2) == false", false},
		{"Test 16", "(1 > 2) == true", false},
		{"Test 17", "(1 > 2) == false", true},
		{"Test 18", "!true", false},
		{"Test 19", "!false", true},
		{"Test 20", "!5", false},
		{"Test 21", "!!true", true},
		{"Test 22", "!!false", false},
		{"Test 23", "!!5", true},
	}

	runVmTests(t, testCases)
}

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
			"object has wrong value. want=%d, got =%d",
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
	}

	runVmTests(t, testCases)
}

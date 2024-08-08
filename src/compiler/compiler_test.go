package compiler

import (
	"fmt"
	"testing"

	"github.com/tjapit/monkey/src/ast"
	"github.com/tjapit/monkey/src/code"
	"github.com/tjapit/monkey/src/lexer"
	"github.com/tjapit/monkey/src/object"
	"github.com/tjapit/monkey/src/parser"
)

type compilerTestCase struct {
	desc                 string
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}
	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}

func testInstructions(
	expected []code.Instructions,
	actual code.Instructions,
) error {
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf(
			"wrong instructions length.\nwant=%q\ngot =%q",
			concatted,
			actual,
		)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf(
				"wrong instruction at %d.\nwant=%q\ngot =%q",
				i,
				concatted,
				actual,
			)
		}
	}

	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
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

func testConstants(
	t *testing.T,
	expected []interface{},
	actual []object.Object,
) error {
	if len(expected) != len(actual) {
		return fmt.Errorf(
			"wrong number of constants. want=%d, got =%d",
			len(expected),
			len(actual),
		)
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		}
	}

	return nil
}

func runCompilerTests(t *testing.T, testCases []compilerTestCase) {
	t.Helper()

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			program := parse(tC.input)

			compiler := New()
			err := compiler.Compile(program)
			if err != nil {
				t.Fatalf("compiler error: %s", err)
			}

			bytecode := compiler.Bytecode()

			err = testInstructions(
				tC.expectedInstructions,
				bytecode.Instructions,
			)
			if err != nil {
				t.Fatalf("testInstructions failed: %s", err)
			}

			err = testConstants(t, tC.expectedConstants, bytecode.Constants)
			if err != nil {
				t.Fatalf("testConstants failed: %s", err)
			}
		})
	}
}

func TestIntegerArithmetic(t *testing.T) {
	testCases := []compilerTestCase{
		{
			desc:              "Test 1",
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "Test 2",
			input:             "1; 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "Test 3",
			input:             "1 - 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "Test 4",
			input:             "1 * 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "Test 5",
			input:             "2 / 1",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "Test 6",
			input:             "-1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testCases)
}

func TestBooleanExpression(t *testing.T) {
	testCases := []compilerTestCase{
		{
			desc:              "Test 1",
			input:             "true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "Test 2",
			input:             "false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "Test 3",
			input:             "1 > 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "Test 4",
			input:             "1 < 2",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "Test 5",
			input:             "1 == 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "Test 6",
			input:             "1 != 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "Test 7",
			input:             "true == false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "Test 8",
			input:             "true != false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "Test 9",
			input:             "!true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testCases)
}

func TestConditionals(t *testing.T) {
	testCases := []compilerTestCase{
		{
			desc: "Test 1",
			input: `
      if (true) { 10 }; 3333;
      `,
			expectedConstants: []interface{}{10, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 7),
				// 0004
				code.Make(code.OpConstant, 0),
				// 0007
				code.Make(code.OpPop),
				// 0008
				code.Make(code.OpConstant, 1),
				// 0011
				code.Make(code.OpPop),
			},
		},
		{
			desc: "Test 2",
			input: `
      if (true) { 10 } else { 20 }; 3333;
      `,
			expectedConstants: []interface{}{10, 20, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 10),
				// 0004 | first line of Consequence
				code.Make(code.OpConstant, 0),
				// 0007
				code.Make(code.OpJump, 13),
				// 0010 | first line of Alternative
				code.Make(code.OpConstant, 1),
				// 0013
				code.Make(code.OpPop),
				// 0014
				code.Make(code.OpConstant, 2),
				// 0017
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testCases)
}

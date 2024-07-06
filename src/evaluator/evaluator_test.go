package evaluator

import (
	"fmt"
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
	env := object.NewEnvironment()

	return Eval(program, env)
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

func TestEvalBooleanExpression(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected bool
	}{
		{"Test 1", "true", true},
		{"Test 2", "false", false},
		{"Test 3", "1 < 2", true},
		{"Test 4", "1 > 2", false},
		{"Test 5", "1 < 1", false},
		{"Test 6", "1 > 1", false},
		{"Test 7", "1 == 1", true},
		{"Test 8", "1 != 1", false},
		{"Test 9", "1 == 2", false},
		{"Test 10", "1 != 2", true},
		{"Test 11", "true == true", true},
		{"Test 12", "false == false", true},
		{"Test 13", "true == false", false},
		{"Test 14", "true != false", true},
		{"Test 15", "false != true", true},
		{"Test 16", "(1 < 2) == true", true},
		{"Test 17", "(1 < 2) == false", false},
		{"Test 18", "(1 > 2) == true", false},
		{"Test 19", "(1 > 2) == false", true},
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

func TestIfElseExpressions(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected interface{}
	}{
		{"Test 1", "if (true) { 10 }", 10},
		{"Test 2", "if (false) { 10 }", nil},
		{"Test 3", "if (1) { 10 }", 10},
		{"Test 4", "if (1 < 2) { 10 }", 10},
		{"Test 5", "if (1 > 2) { 10 }", nil},
		{"Test 6", "if (1 > 2) { 10 } else { 20 }", 20},
		{"Test 7", "if (1 < 2) { 10 } else { 20 }", 10},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			evaluated := testEval(tC.input)
			integer, ok := tC.expected.(int)
			if ok {
				testIntegerObject(t, evaluated, int64(integer))
			} else {
				testNullObject(t, evaluated)
			}
		})
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected int64
	}{
		{"Test 1", "return 10;", 10},
		{"Test 2", "return 10; 9;", 10},
		{"Test 3", "return 2 * 5; 9;", 10},
		{"Test 4", "9; return 2 * 5; 9;", 10},
		{"Test 5", `
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }

  return 1;
}`, 10},
		{"Test 6", `
if (10 > 1) {
  if (10 < 1) {
    return 10;
  }

  return 1;
}`, 1},
		{"Test 7", `
let f = fn(x) {
  return x;
  x + 10;
};
f(10);`, 10},
		{"Test 7", `
let f = fn(x) {
  let result = x + 10;
  return result;
  return 10;
};
f(10);`, 20},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			evaluated := testEval(tC.input)
			testIntegerObject(t, evaluated, tC.expected)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "Test 1",
			input:    "5 + true;",
			expected: "type mismatch: INTEGER + BOOLEAN",
		},
		{
			desc:     "Test 2",
			input:    "5 + true; 5;",
			expected: "type mismatch: INTEGER + BOOLEAN",
		},
		{
			desc:     "Test 3",
			input:    "-true",
			expected: "unknown operator: -BOOLEAN",
		},
		{
			desc:     "Test 4",
			input:    "true + false;",
			expected: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			desc:     "Test 5",
			input:    "5; true + false; 5",
			expected: "unknown operator: BOOLEAN + BOOLEAN",
		},

		{
			desc:     "Test 6",
			input:    "if (10 > 1) { true + false; }",
			expected: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			desc: "Test 7",
			input: `
      if (10 > 1) {
        if (10 > 1) {
          return true + false;
        }

        return 1;
      }`,
			expected: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			desc:     "Test 8",
			input:    "foobar",
			expected: "identifier not found: foobar",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			evaluated := testEval(tC.input)

			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf(
					"no error object returned. got=%T(%+v)",
					evaluated,
					evaluated,
				)
				t.FailNow()
			}
			if errObj.Message != tC.expected {
				t.Errorf(
					"wrong error message. want=%q, got=%q",
					tC.expected,
					errObj.Message,
				)
			}
		})
	}
}

func TestLetStatements(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected int64
	}{
		{
			desc:     "Test 1",
			input:    "let a = 5; a;",
			expected: 5,
		},
		{
			desc:     "Test 2",
			input:    "let a = 5 * 5; a;",
			expected: 25,
		},
		{
			desc:     "Test 3",
			input:    "let a = 5; let b = a; b;",
			expected: 5,
		},
		{
			desc:     "Test 4",
			input:    "let a = 5; let b = a; let c = a + b + 5; c;",
			expected: 15,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			testIntegerObject(t, testEval(tC.input), tC.expected)
		})
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameters is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected int64
	}{
		{
			desc:     "Test 1",
			input:    "let identity = fn(x) { x; }; identity(5);",
			expected: 5,
		},
		{
			desc:     "Test 2",
			input:    "let identity = fn(x) { return x; }; identity(5);",
			expected: 5,
		},
		{
			desc:     "Test 3",
			input:    "let double = fn(x) { x * 2; }; double(5);",
			expected: 10,
		},
		{
			desc:     "Test 4",
			input:    "let add = fn(x, y) { x + y; }; add(5, 5);",
			expected: 10,
		},
		{
			desc:     "Test 5",
			input:    "let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));",
			expected: 20,
		},
		{
			desc:     "Test 6",
			input:    "fn(x) { x; }(5)",
			expected: 5,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			testIntegerObject(t, testEval(tC.input), tC.expected)
		})
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
  fn(y) { x + y; };
};

let addTwo = newAdder(2);
addTwo(2);`

	expected := 4
	testIntegerObject(t, testEval(input), int64(expected))
}

func TestStringLiteral(t *testing.T) {
	expected := "Hello World!"
	input := fmt.Sprintf(`"%s"`, expected)

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != expected {
		t.Errorf("String has wrong value. want=%q, got=%q", expected, str.Value)
	}
}

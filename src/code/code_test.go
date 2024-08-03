package code

import (
	"testing"
)

func TestMake(t *testing.T) {
	testCases := []struct {
		desc     string
		op       Opcode
		operands []int
		expected []byte
	}{
		{
			desc:     "Test 1",
			op:       OpConstant,
			operands: []int{65534},
			expected: []byte{byte(OpConstant), 255, 254},
		},
		{"Test 2", OpAdd, []int{}, []byte{byte(OpAdd)}},
		{"Test 3", OpPop, []int{}, []byte{byte(OpPop)}},
		{"Test 4", OpSub, []int{}, []byte{byte(OpSub)}},
		{"Test 5", OpMul, []int{}, []byte{byte(OpMul)}},
		{"Test 6", OpDiv, []int{}, []byte{byte(OpDiv)}},
		{"Test 7", OpTrue, []int{}, []byte{byte(OpTrue)}},
		{"Test 8", OpFalse, []int{}, []byte{byte(OpFalse)}},
		{"Test 9", OpEqual, []int{}, []byte{byte(OpEqual)}},
		{"Test 10", OpNotEqual, []int{}, []byte{byte(OpNotEqual)}},
		{"Test 11", OpGreaterThan, []int{}, []byte{byte(OpGreaterThan)}},
		{"Test 12", OpMinus, []int{}, []byte{byte(OpMinus)}},
		{"Test 13", OpBang, []int{}, []byte{byte(OpBang)}},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			instruction := Make(tC.op, tC.operands...)

			if len(instruction) != len(tC.expected) {
				t.Errorf(
					"instruction has wrong length. want=%d, got =%d",
					len(tC.expected),
					len(instruction),
				)
			}

			for i, b := range tC.expected {
				if instruction[i] != b {
					t.Errorf(
						"wrong byte at pos %d. want=%q, got =%q",
						i,
						b,
						instruction[i],
					)
				}
			}
		})
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpPop),
		Make(OpSub),
		Make(OpMul),
		Make(OpDiv),
		Make(OpTrue),
		Make(OpFalse),
		Make(OpEqual),
		Make(OpNotEqual),
		Make(OpGreaterThan),
		Make(OpMinus),
		Make(OpBang),
	}

	expected := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
0007 OpPop
0008 OpSub
0009 OpMul
0010 OpDiv
0011 OpTrue
0012 OpFalse
0013 OpEqual
0014 OpNotEqual
0015 OpGreaterThan
0016 OpMinus
0017 OpBang
`

	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf(
			"instructions wrongly formatted.\nwant=%q\ngot =%q",
			expected,
			concatted.String(),
		)
	}
}

func TestReadOperands(t *testing.T) {
	testCases := []struct {
		desc      string
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{"Test 1", OpConstant, []int{65535}, 2},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			instruction := Make(tC.op, tC.operands...)

			def, err := Lookup(byte(tC.op))
			if err != nil {
				t.Fatalf("definition not found: %q", err)
			}

			operandsRead, nBytes := ReadOperands(def, instruction[1:])
			if nBytes != tC.bytesRead {
				t.Fatalf("nBytes wrong. want=%d, got =%d", tC.bytesRead, nBytes)
			}

			for i, want := range tC.operands {
				if operandsRead[i] != want {
					t.Errorf(
						"operand wrong. want=%d, got =%d",
						want,
						operandsRead[i],
					)
				}
			}
		})
	}
}

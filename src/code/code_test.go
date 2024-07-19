package code

import "testing"

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
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			instruction := Make(tC.op, tC.operands...)

			if len(instruction) != len(tC.expected) {
				t.Errorf(
					"instruction has wrong length. want=%d, got=%d",
					len(tC.expected),
					len(instruction),
				)
			}

			for i, b := range tC.expected {
				if instruction[i] != b {
					t.Errorf(
						"wrong byte at pos %d. want=%d, got=%d",
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
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}

	expected := `0000 OpConstant 1
0003 OpConstant 2
0006 OpConstant 65535
`

	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf(
			"instructions wrongly formatted.\nwant=%q\ngot=%q",
			expected,
			concatted.String(),
		)
	}
}

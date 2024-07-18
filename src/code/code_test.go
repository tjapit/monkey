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

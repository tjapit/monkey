package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota
)

type Definiton struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definiton{
	OpConstant: {"OpConstant", []int{2}},
}

func Lookup(op byte) (*Definiton, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

// Encodes the given Opcode and operands into a bytecode.
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, width := range def.OperandWidths {
		instructionLen += width
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, operand := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		}
		offset += width
	}

	return instruction
}

func (ins *Instructions) String() string {
	return ""
}

// Decodes bytecode made by Make(). Returns operands and number of bytes read.
func ReadOperands(def *Definiton, operandsBytes []byte) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(operandsBytes[offset:]))
		}

		offset += width
	}
	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

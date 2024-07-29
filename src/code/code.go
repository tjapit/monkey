package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
)

type Definiton struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definiton{
	OpConstant: {"OpConstant", []int{2}},
	OpAdd:      {"OpAdd", []int{}},
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

// Mini-disassembler
func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read // the 1 is for the Opcode (i.e. Opcode + num of bytes read)
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definiton, operands []int) string {
	operandCount := len(def.OperandWidths)
	if len(operands) != operandCount {
		return fmt.Sprintf(
			"ERROR: operand len %d does not match defined %d\n",
			len(operands),
			operandCount,
		)
	}

	switch operandCount {
	case 0:
		return fmt.Sprintf("%s", def.Name)
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
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

func ReadUint16(operand Instructions) uint16 {
	return binary.BigEndian.Uint16(operand)
}
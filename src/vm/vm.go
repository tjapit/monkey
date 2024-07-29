package vm

import (
	"github.com/tjapit/monkey/src/compiler"
	"github.com/tjapit/monkey/src/object"
)

type VM struct{}

func New(bytecode *compiler.Bytecode) VM {
	return VM{}
}

func (vm *VM) Run() error {
	return nil
}

func (vm *VM) StackTop() object.Object {
	return nil
}

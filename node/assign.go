package node

import (
	"errors"
	"fmt"

	"jakobvarmose/compute/vtype"
)

type Assign struct {
	left  *Variable
	right Node
}

func NewAssign(e Environment, left, right Node) *Assign {
	return &Assign{left.(*Variable), right}
}

func (op *Assign) Name() string {
	return "="
}
func (op *Assign) Left() Node {
	return op.left
}
func (op *Assign) Right() Node {
	return op.right
}

func (op *Assign) String() string {
	return binaryOpString(op)
}

func (op *Assign) MarshalString() string {
	return binaryOpString(op)
}

func (op *Assign) Eval(e Environment) Node {
	left := op.left
	right := op.right.Eval(e)
	if left.Name == "@" {
		switch right := right.(type) {
		case (*String):
			fmt.Println(right.Value)
		default:
			fmt.Println(right)
		}
		return right
	} else {
		e.Set(left.Name, right)
	}
	return right
}

func (op *Assign) Convert(e Environment, t Type) (Node, error) {
	switch t {
	case TypeUnknown:
		return op, nil
	default:
		return nil, errors.New("cannot convert")
	}
}

func (op *Assign) Type() Type {
	return TypeUnknown
}

func (op *Assign) ComputedType(e Environment) vtype.VType {
	return op.right.ComputedType(e)
}

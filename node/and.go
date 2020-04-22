package node

import (
	"errors"
	"jakobvarmose/compute/vtype"
)

type And struct {
	left, right Node
}

func NewAnd(e Environment, left, right Node) *And {
	if left.ComputedType(e) != vtype.Boolean {
		panic("invalid type for and operator")
	}
	if right.ComputedType(e) != vtype.Boolean {
		panic("invalid type for and operator")
	}
	return &And{left, right}
}

func (op *And) Name() string {
	return "and"
}
func (op *And) Left() Node {
	return op.left
}
func (op *And) Right() Node {
	return op.right
}

func (op *And) String() string {
	return binaryOpString(op)
}

func (op *And) MarshalString() string {
	return binaryOpString(op)
}

func (op *And) Eval(e Environment) Node {
	left := op.left.Eval(e)
	right := op.right.Eval(e)
	if boolLeft, ok := left.(*Bool); ok {
		if boolLeft.Val {
			return right
		}
		return left
	}
	if boolRight, ok := right.(*Bool); ok {
		if boolRight.Val {
			return left
		}
		return right
	}
	if andRight, ok := right.(*And); ok {
		return (&And{
			&And{left, andRight.left},
			andRight.right,
		}).Eval(e)
	}
	return &And{left, right}
}

func (op *And) Convert(e Environment, t Type) (Node, error) {
	switch t {
	case TypeUnknown:
		return op, nil
	default:
		return nil, errors.New("cannot convert")
	}
}

func (op *And) Type() Type {
	return TypeUnknown
}

func (op *And) ComputedType(e Environment) vtype.VType {
	return vtype.Boolean
}

package node

import (
	"errors"
	"jakobvarmose/compute/vtype"
)

type Or struct {
	left, right Node
}

func NewOr(e Environment, left, right Node) *Or {
	if left.ComputedType(e) != vtype.Boolean {
		panic("invalid type for or operator")
	}
	if right.ComputedType(e) != vtype.Boolean {
		panic("invalid type for or operator")
	}
	return &Or{left, right}
}

func (op *Or) Name() string {
	return "or"
}
func (op *Or) Left() Node {
	return op.left
}
func (op *Or) Right() Node {
	return op.right
}

func (op *Or) String() string {
	return binaryOpString(op)
}

func (op *Or) MarshalString() string {
	return binaryOpString(op)
}

func (op *Or) Eval(e Environment) Node {
	left := op.left.Eval(e)
	right := op.right.Eval(e)
	if boolLeft, ok := left.(*Bool); ok {
		if boolLeft.Val {
			return left
		}
		return right
	}
	if boolRight, ok := right.(*Bool); ok {
		if boolRight.Val {
			return right
		}
		return left
	}
	if orRight, ok := right.(*Or); ok {
		return (&Or{
			&Or{left, orRight.left},
			orRight.right,
		}).Eval(e)
	}
	return &Or{left, right}
}

func (op *Or) Convert(e Environment, t Type) (Node, error) {
	switch t {
	case TypeUnknown:
		return op, nil
	default:
		return nil, errors.New("cannot convert")
	}
}

func (op *Or) Type() Type {
	return TypeUnknown
}

func (op *Or) ComputedType(e Environment) vtype.VType {
	return vtype.Boolean
}

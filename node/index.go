package node

import (
	"errors"
	"jakobvarmose/compute/vtype"
	"math/big"
)

type Index struct {
	left, right Node
}

func NewIndex(e Environment, left, right Node) *Index {
	if _, ok := left.ComputedType(e).(*vtype.ListType); !ok {
		panic("wrong left type for indexing")
	}
	return &Index{left, right}
}

func (op *Index) Name() string {
	return "#"
}
func (op *Index) Left() Node {
	return op.left
}
func (op *Index) Right() Node {
	return op.right
}

func (op *Index) String() string {
	return binaryOpString(op)
}

func (op *Index) MarshalString() string {
	return binaryOpString(op)
}

func (op *Index) Eval(e Environment) Node {
	left := op.left.Eval(e)
	right := op.right.Eval(e)
	if left, ok := left.(ListInterface); ok {
		right, err := right.Convert(e, TypeInt)
		if err != nil {
			panic(err)
		}
		if right.(*Int).Big.Cmp(big.NewInt(0)) < 0 ||
			right.(*Int).Big.Cmp(left.Len()) >= 0 {
			panic("index out of range")
		}
		return left.ItemAt(right.(*Int).Big).Eval(e)
	}
	return &Index{left, right}
}

func (op *Index) Convert(e Environment, t Type) (Node, error) {
	switch t {
	case TypeUnknown:
		return op, nil
	default:
		return nil, errors.New("cannot convert")
	}
}

func (op *Index) Type() Type {
	return TypeUnknown
}

func (op *Index) ComputedType(e Environment) vtype.VType {
	return op.left.ComputedType(e).(*vtype.ListType).Item
}

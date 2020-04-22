package node

import (
	"errors"
	"jakobvarmose/compute/vtype"
	"math/big"
)

type Minus struct {
	left, right Node
}

func NewMinus(left, right Node) *Minus {
	return &Minus{left, right}
}

func (op *Minus) Name() string {
	return "-"
}
func (op *Minus) Left() Node {
	return op.left
}
func (op *Minus) Right() Node {
	return op.right
}

func (op *Minus) String() string {
	return binaryOpString(op)
}

func (op *Minus) MarshalString() string {
	return binaryOpString(op)
}

func (op *Minus) Eval(e Environment) Node {
	left := op.left.Eval(e)
	right := op.right.Eval(e)
	leftType := left.ComputedType(e)
	rightType := right.ComputedType(e)
	if leftType == vtype.Complex || rightType == vtype.Complex {
		c1, _ := left.Convert(e, TypeComplex)
		c2, _ := right.Convert(e, TypeComplex)
		return &Complex{
			new(big.Rat).Sub(c1.(*Complex).Real, c2.(*Complex).Real),
			new(big.Rat).Sub(c1.(*Complex).Imag, c2.(*Complex).Imag),
		}
	}
	if leftType == vtype.Rational || rightType == vtype.Rational {
		c1, _ := left.Convert(e, TypeRat)
		c2, _ := right.Convert(e, TypeRat)
		return &Rat{
			new(big.Rat).Sub(c1.(*Rat).Big, c2.(*Rat).Big),
		}
	}
	if leftType == vtype.Integer || rightType == vtype.Integer {
		c1, _ := left.Convert(e, TypeInt)
		c2, _ := right.Convert(e, TypeInt)
		return &Int{
			new(big.Int).Sub(c1.(*Int).Big, c2.(*Int).Big),
		}
	}
	return &Minus{left, right}
}

func (op *Minus) Convert(e Environment, t Type) (Node, error) {
	switch t {
	case TypeUnknown:
		return op, nil
	default:
		return nil, errors.New("cannot convert")
	}
}

func (op *Minus) Type() Type {
	return TypeUnknown
}

func (op *Minus) ComputedType(e Environment) vtype.VType {
	leftType := op.left.ComputedType(e)
	rightType := op.right.ComputedType(e)
	if leftType == vtype.Complex || rightType == vtype.Complex {
		return vtype.Complex
	}
	if leftType == vtype.Real || rightType == vtype.Real {
		return vtype.Real
	}
	if leftType == vtype.Rational || rightType == vtype.Rational {
		return vtype.Rational
	}
	if leftType == vtype.Integer || rightType == vtype.Integer {
		return vtype.Integer
	}
	return vtype.Boolean
}

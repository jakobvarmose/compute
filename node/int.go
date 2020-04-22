package node

import (
	"errors"
	"math/big"

	"jakobvarmose/compute/vtype"

	"github.com/ericlagergren/decimal"
)

type Int struct {
	Big *big.Int
}

func NewInt(str string) (*Int, error) {
	big, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return nil, errors.New("invalid integer")
	}
	return &Int{big}, nil
}

func (i *Int) String() string {
	return i.Big.String()
}
func (i *Int) MarshalString() string {
	return i.String()
}

func (i *Int) Eval(e Environment) Node {
	return i
}

func (i *Int) Convert(e Environment, t Type) (Node, error) {
	switch t {
	case TypeInt:
		return i, nil
	case TypeRat:
		return &Rat{new(big.Rat).SetInt(i.Big)}, nil
	case TypeDecimal:
		return &Decimal{decimal.WithPrecision(e.Precision()).SetBigMantScale(i.Big, 0)}, nil
	case TypeComplex:
		return &Complex{new(big.Rat).SetInt(i.Big), new(big.Rat)}, nil
	case TypeString:
		return &String{i.String()}, nil
	}
	return nil, errors.New("cannot convert")
}

func (i *Int) Type() Type {
	return TypeInt
}

func (i *Int) ComputedType(e Environment) vtype.VType {
	return vtype.Integer
}

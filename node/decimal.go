package node

import (
	"errors"
	"fmt"
	"math/big"

	"jakobvarmose/compute/vtype"

	"github.com/ericlagergren/decimal"
)

type Decimal struct {
	Val *decimal.Big
}

func NewDecimal(e Environment, str string) (*Decimal, error) {
	val, ok := decimal.WithPrecision(e.Precision()).SetString(str)
	if !ok {
		return nil, errors.New("invalid decimal")
	}
	return &Decimal{val}, nil
}

func (d *Decimal) String() string {
	return fmt.Sprintf("%.6e", d.Val)
}
func (d *Decimal) MarshalString() string {
	buf, _ := d.Val.MarshalText()
	return string(buf)
}

func (d *Decimal) Eval(e Environment) Node {
	return d
}

func (d *Decimal) Convert(env Environment, t Type) (Node, error) {
	switch t {
	case TypeInt:
		if d.Val.IsInt() {
			return &Int{d.Val.Int(nil)}, nil
		}
	case TypeRat:
		return &Rat{d.Val.Rat(nil)}, nil
	case TypeDecimal:
		return d, nil
	case TypeComplex:
		return &Complex{d.Val.Rat(nil), new(big.Rat)}, nil
	case TypeString:
		return &String{d.String()}, nil
	}
	return nil, errors.New("cannot convert")
}

func (d *Decimal) Type() Type {
	return TypeDecimal
}

func (d *Decimal) ComputedType(e Environment) vtype.VType {
	return vtype.Real
}

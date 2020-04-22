package node

import (
	"errors"
	"math/big"

	"jakobvarmose/compute/vtype"

	"github.com/ericlagergren/decimal"
)

type Rat struct {
	Big *big.Rat
}

func NewRat(str string) (*Rat, error) {
	big, ok := new(big.Rat).SetString(str)
	if !ok {
		return nil, errors.New("invalid rat")
	}
	return &Rat{big}, nil
}

func (r *Rat) String() string {
	return r.Big.RatString()
}
func (r *Rat) MarshalString() string {
	return r.Big.RatString()
}

func (r *Rat) Eval(e Environment) Node {
	return r
}

func (r *Rat) Convert(e Environment, t Type) (Node, error) {
	switch t {
	case TypeInt:
		if r.Big.IsInt() {
			return &Int{r.Big.Num()}, nil
		}
	case TypeRat:
		return r, nil
	case TypeDecimal:
		return &Decimal{decimal.WithPrecision(e.Precision()).SetRat(r.Big)}, nil
	case TypeComplex:
		return &Complex{r.Big, new(big.Rat)}, nil
	case TypeString:
		return &String{r.String()}, nil
	}
	return nil, errors.New("cannot convert")
}

func (r *Rat) Type() Type {
	return TypeRat
}

func (r *Rat) ComputedType(e Environment) vtype.VType {
	return vtype.Rational
}

func (r *Rat) Floor() *Int {
	return &Int{new(big.Int).Div(r.Big.Num(), r.Big.Denom())}
}

func (r *Rat) Ceil() *Int {
	i := new(big.Int).Add(r.Big.Num(), r.Big.Denom())
	i.Sub(i, big.NewInt(1))
	return &Int{i.Div(i, r.Big.Denom())}
}

package node

import (
	"errors"
	"jakobvarmose/compute/vtype"

	"github.com/ericlagergren/decimal"
	dmath "github.com/ericlagergren/decimal/math"
)

var Exacts = []*Exact{
	{"E"},
	{"PI"},
}

type Exact struct {
	Name string
}

func NewExact(str string) (*Exact, error) {
	return &Exact{str}, nil
}

func (e *Exact) String() string {
	return e.Name
}
func (e *Exact) MarshalString() string {
	return e.Name
}

func (e *Exact) Eval(env Environment) Node {
	return e
}

func (e *Exact) Convert(env Environment, t Type) (Node, error) {
	switch t {
	case TypeDecimal:
		d := decimal.WithPrecision(env.Precision())
		switch e.Name {
		case "E":
			return &Decimal{dmath.E(d)}, nil
		case "PI":
			return &Decimal{dmath.Pi(d)}, nil
		}
	case TypeExact:
		return e, nil
	case TypeString:
		return &String{e.String()}, nil
	}
	return nil, errors.New("cannot convert")
}

func (e *Exact) Type() Type {
	return TypeExact
}

func (e *Exact) ComputedType(env Environment) vtype.VType {
	return vtype.Real
}

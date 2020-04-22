package node

import (
	"errors"
	"fmt"

	"jakobvarmose/compute/vtype"
)

type Bool struct {
	Val bool
}

func NewBool(val string) (*Bool, error) {
	switch val {
	case "true":
		return &Bool{true}, nil
	case "false":
		return &Bool{false}, nil
	default:
		return nil, errors.New("cannot parse bool")
	}
}

func (b *Bool) String() string {
	return fmt.Sprintf("%t", b.Val)
}
func (b *Bool) MarshalString() string {
	return b.String()
}
func (b *Bool) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%t", b.Val)
}

func (b *Bool) Eval(e Environment) Node {
	return b
}

func (b *Bool) Convert(e Environment, t Type) (Node, error) {
	switch t {
	case TypeBool:
		return b, nil
	case TypeString:
		return &String{b.String()}, nil
	}
	return nil, errors.New("cannot convert bool")
}

func (b *Bool) Type() Type {
	return TypeBool
}

func (b *Bool) ComputedType(e Environment) vtype.VType {
	return vtype.Boolean
}

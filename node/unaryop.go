package node

import (
	"errors"
	"jakobvarmose/compute/vtype"
)

type UnaryOp struct {
	arg Node
	str string
}

func NewUnaryOp(arg Node, str string) *UnaryOp {
	return &UnaryOp{arg, str}
}

func (u *UnaryOp) String() string {
	argStr := u.arg.String()
	if _, ok := u.arg.(*BinaryOp); ok || u.str[0] >= 'a' && u.str[0] <= 'z' || u.str[0] == '@' {
		argStr = "(" + argStr + ")"
	}
	return u.str + argStr
}
func (u *UnaryOp) MarshalString() string {
	return u.str + "(" + u.arg.String() + ")"
}

func (u *UnaryOp) Eval(e Environment) Node {
	return unary_call(e, u)
}

func (u *UnaryOp) Convert(e Environment, t Type) (Node, error) {
	switch t {
	default:
		return nil, errors.New("cannot convert")
	}
}

func (u *UnaryOp) Type() Type {
	return TypeUnknown
}

func (b *UnaryOp) ComputedType(e Environment) vtype.VType {
	return vtype.Boolean
}

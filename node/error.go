package node

import (
	"errors"

	"jakobvarmose/compute/vtype"
)

type Error struct {
	text string
	vt   vtype.VType
}

func NewError(text string) *Error {
	return &Error{text, vtype.Boolean}
}

func (f *Error) String() string {
	return f.text
}
func (f *Error) MarshalString() string {
	return "error(\"" + f.text + "\")"
}

func (f *Error) Eval(e Environment) Node {
	return f
}

func (f *Error) Convert(env Environment, t Type) (Node, error) {
	switch t {
	case TypeError:
		return f, nil
	default:
		return nil, errors.New("cannot convert")
	}
}

func (f *Error) Type() Type {
	return TypeError
}

func (f *Error) ComputedType(e Environment) vtype.VType {
	return f.vt
}

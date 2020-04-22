package node

import (
	"errors"

	"jakobvarmose/compute/vtype"
)

type Builtin struct {
	Name string
	Fn   func(e Environment, arg Node, str string) Node
	T    func(e Environment, arg Node, str string) vtype.VType
	Typ  *vtype.ConvertType
}

func (b *Builtin) String() string {
	return b.Name
}
func (b *Builtin) MarshalString() string {
	return b.String()
}

func (b *Builtin) Eval(e Environment) Node {
	return b
}

func (b *Builtin) Convert(env Environment, t Type) (Node, error) {
	switch t {
	case TypeBuiltin:
		return b, nil
	default:
		return nil, errors.New("cannot convert")
	}
}

func (b *Builtin) Type() Type {
	return TypeBuiltin
}

func (b *Builtin) ComputedType(e Environment) vtype.VType {
	return b.Typ
}

package node

import (
	"jakobvarmose/compute/vtype"
)

type Environment interface {
	Get(name string) Node
	Set(name string, value Node)

	GetType(name string) vtype.VType
	SetType(name string, vt vtype.VType)

	SideEffects() bool
	SetSideEffects(val bool)

	SetPrecision(int)
	Precision() int
}

type Node interface {
	String() string
	MarshalString() string
	Eval(env Environment) Node
	Convert(env Environment, t Type) (Node, error)
	Type() Type

	ComputedType(e Environment) vtype.VType
}

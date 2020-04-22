package node

import (
	"errors"
	"fmt"

	"jakobvarmose/compute/vtype"
)

var (
	variable_counter = 1
)

type Variable struct {
	Name string
	typ  vtype.VType
}

func NewVariable(name string) *Variable {
	return &Variable{name, nil}
}

func UniqueVariable(typ vtype.VType) *Variable {
	variable_counter++
	return &Variable{fmt.Sprintf("@%d", variable_counter), typ}
}
func (v *Variable) String() string {
	return v.Name
}
func (v *Variable) MarshalString() string {
	return v.Name
}

func (v *Variable) Eval(e Environment) Node {
	value := e.Get(v.Name)
	if value == nil {
		return v
	}
	return value.Eval(e)
}

func (v *Variable) Convert(env Environment, t Type) (Node, error) {
	switch t {
	case TypeVariable:
		return v, nil
	default:
		return nil, errors.New("cannot convert")
	}
}

func (v *Variable) Type() Type {
	return TypeVariable
}

func (v *Variable) ComputedType(e Environment) vtype.VType {
	if v.typ != nil {
		return v.typ
	}
	return e.GetType(v.Name)
}

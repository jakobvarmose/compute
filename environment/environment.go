package environment

import (
	"encoding/json"
	"io/ioutil"

	"jakobvarmose/compute/node"
	"jakobvarmose/compute/vtype"
)

type Environment interface {
	Get(name string) node.Node
	Set(name string, value node.Node)
	SetDontSave(name string, value node.Node)

	SetType(name string, vt vtype.VType)
	GetType(name string) vtype.VType

	Save() error

	SideEffects() bool
	SetSideEffects(val bool)

	SetPrecision(int)
	Precision() int
}

type environment struct {
	variables   map[string]variable
	types       map[string]vtype.VType
	sideEffects bool
	precision   int
}

type variable struct {
	val  node.Node
	save bool
}

func New() Environment {
	return &environment{
		variables:   make(map[string]variable),
		types:       make(map[string]vtype.VType),
		sideEffects: false,
		precision:   16,
	}
}

func (e *environment) Get(name string) node.Node {
	return e.variables[name].val
}

func (e *environment) Set(name string, value node.Node) {
	delete(e.types, name)
	if value == nil {
		delete(e.variables, name)
	} else {
		e.variables[name] = variable{
			val:  value,
			save: true,
		}
	}
}

func (e *environment) SetDontSave(name string, value node.Node) {
	delete(e.types, name)
	if value == nil {
		delete(e.variables, name)
	} else {
		e.variables[name] = variable{
			val:  value,
			save: false,
		}
	}
}

func (e *environment) SetType(name string, vt vtype.VType) {
	delete(e.variables, name)
	e.types[name] = vt
}

func (e *environment) GetType(name string) vtype.VType {
	if typ, ok := e.types[name]; ok {
		return typ
	}
	if v, ok := e.variables[name]; ok {
		return v.val.ComputedType(e)
	}
	return nil
}

/*func (e *environment) Load() error {
	buf, err := ioutil.ReadFile("env.json")
	if err != nil {
		return err
	}
	var variables map[string]string
	err = json.Unmarshal(buf, variables)
	if err != nil {
		return err
	}
	for name, value := range variables {
		e.variables[name] = variable{
			val:  value, // TODO parse
			save: true,
		}
	}
	return nil
}*/

func (e *environment) Save() error {
	variables := make(map[string]string)
	for name, value := range e.variables {
		if value.save {
			variables[name] = value.val.MarshalString()
		}
	}
	buf, err := json.MarshalIndent(map[string]interface{}{
		"variables": variables,
		"types":     e.types,
	}, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("env.json", buf, 0700)
	if err != nil {
		return err
	}
	return nil
}
func (e *environment) SideEffects() bool {
	return e.sideEffects
}
func (e *environment) SetSideEffects(val bool) {
	e.sideEffects = val
}

func (e *environment) SetPrecision(prec int) {
	e.precision = prec
}

func (e *environment) Precision() int {
	return e.precision
}

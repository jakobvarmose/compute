package unary

import (
	"jakobvarmose/compute/node"
	"jakobvarmose/compute/vtype"
)

func CreateUnary(name string,
	calc func(e node.Environment, arg node.Node, str string) node.Node,
) {
	Builtins[name] = &node.Builtin{name, calc,
		func(e node.Environment, arg node.Node, str string) vtype.VType {
			return vtype.Boolean
		},
		&vtype.ConvertType{vtype.Boolean, vtype.Boolean},
	}
}

func CreateUnary2(
	name string,
	calc func(e node.Environment, arg node.Node, str string) node.Node,
	t func(e node.Environment, arg node.Node, str string) vtype.VType,
	typ *vtype.ConvertType,
) {
	Builtins[name] = &node.Builtin{name, calc, t, typ}
}

var Builtins = make(map[string]*node.Builtin)

package binary

import (
	"errors"
	"jakobvarmose/compute/node"
	"jakobvarmose/compute/vtype"
)

type Operator2 struct {
	Name  string
	Check func(e node.Environment, left node.Node, right node.Node) error
	Eval  func(e node.Environment, left node.Node, right node.Node) node.Node
	Typ   func(e node.Environment, left node.Node, right node.Node) vtype.VType
}

func CreateBinary(
	name string,
	check func(e node.Environment, left node.Node, right node.Node) error,
	eval func(e node.Environment, left node.Node, right node.Node) node.Node,
	typ func(e node.Environment, left node.Node, right node.Node) vtype.VType,
) {
	Builtins[name] = &Operator2{name, check, eval, typ}
	operators2 = append(operators2, &Operator2{name, check, eval, typ})
}

var Builtins = make(map[string]*Operator2)

var operators2 = []*Operator2{
	{
		"and",
		func(e node.Environment, left node.Node, right node.Node) error {
			if left.ComputedType(e) != vtype.Boolean {
				return errors.New("invalid type for and operator")
			}
			if right.ComputedType(e) != vtype.Boolean {
				return errors.New("invalid type for and operator")
			}
			return nil
		},
		func(e node.Environment, left node.Node, right node.Node) node.Node {
			left = left.Eval(e)
			right = right.Eval(e)
			if boolLeft, ok := left.(*node.Bool); ok {
				if boolLeft.Val {
					return right
				}
				return left
			}
			if boolRight, ok := right.(*node.Bool); ok {
				if boolRight.Val {
					return left
				}
				return right
			}
			if andRight, ok := right.(*node.And); ok {
				return node.NewAnd(e,
					node.NewAnd(e, left, andRight.Left()),
					andRight.Right(),
				).Eval(e)
			}
			return node.NewAnd(e, left, right)
		},
		func(e node.Environment, left node.Node, right node.Node) vtype.VType {
			return vtype.Boolean
		},
	},
}

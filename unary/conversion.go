package unary

import "jakobvarmose/compute/node"

func init() {
	CreateUnary(
		"toInt",
		func(e node.Environment, arg node.Node, str string) node.Node {
			res, err := arg.Convert(e, node.TypeInt)
			if err == nil {
				return res
			}
			return node.NewUnaryOp(arg, str)
		},
	)
	CreateUnary(
		"toRational",
		func(e node.Environment, arg node.Node, str string) node.Node {
			res, err := arg.Convert(e, node.TypeRat)
			if err == nil {
				return res
			}
			return node.NewUnaryOp(arg, str)
		},
	)
	CreateUnary(
		"toComplex",
		func(e node.Environment, arg node.Node, str string) node.Node {
			res, err := arg.Convert(e, node.TypeComplex)
			if err == nil {
				return res
			}
			return node.NewUnaryOp(arg, str)
		},
	)
	CreateUnary(
		"toString",
		func(e node.Environment, arg node.Node, str string) node.Node {
			res, err := arg.Convert(e, node.TypeString)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			return res
		},
	)
}

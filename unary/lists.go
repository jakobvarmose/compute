package unary

import (
	"jakobvarmose/compute/node"
	"jakobvarmose/compute/vtype"
	"math/big"
)

func init() {
	// List functions
	CreateUnary2(
		"range",
		func(e node.Environment, arg node.Node, str string) node.Node {
			res, err := arg.Convert(e, node.TypeInt)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			return node.NewInfiniteList(res.(*node.Int).Big)
		},
		func(e node.Environment, arg node.Node, str string) vtype.VType {
			return &vtype.ListType{vtype.Integer}
		},
		&vtype.ConvertType{vtype.Integer, &vtype.ListType{vtype.Integer}},
	)
	CreateUnary(
		"len",
		func(e node.Environment, arg node.Node, str string) node.Node {
			switch arg := arg.(type) {
			case interface {
				Len() *big.Int
			}:
				return &node.Int{arg.Len()}
			}
			return node.NewUnaryOp(arg, str)
		},
	)
	CreateUnary(
		"sum",
		func(e node.Environment, arg node.Node, str string) node.Node {
			if arg, ok := arg.(*node.List); ok {
				return arg.Fold(e, &node.Int{big.NewInt(0)}, "+").Eval(e)
			}
			return node.NewUnaryOp(arg, str)
		},
	)
	CreateUnary(
		"product",
		func(e node.Environment, arg node.Node, str string) node.Node {
			if arg, ok := arg.(*node.List); ok {
				return arg.Fold(e, &node.Int{big.NewInt(1)}, "*").Eval(e)
			}
			return node.NewUnaryOp(arg, str)
		},
	)
	CreateUnary(
		"all", func(e node.Environment, arg node.Node, str string) node.Node {
			if arg, ok := arg.(node.ListInterface); ok {
				unknown := false
				for i := big.NewInt(0); i.Cmp(arg.Len()) < 0; i = new(big.Int).Add(i, big.NewInt(1)) {
					item, err := arg.ItemAt(i).Eval(e).Convert(e, node.TypeBool)
					if err != nil {
						unknown = true
						continue
					}
					if item.(*node.Bool).Val == false {
						return item
					}
				}
				if !unknown {
					return &node.Bool{true}
				}
			}
			return node.NewUnaryOp(arg, str)
		},
	)
	CreateUnary(
		"some",
		func(e node.Environment, arg node.Node, str string) node.Node {
			if arg, ok := arg.(node.ListInterface); ok {
				unknown := false
				for i := big.NewInt(0); i.Cmp(arg.Len()) < 0; i = new(big.Int).Add(i, big.NewInt(1)) {
					item, err := arg.ItemAt(i).Eval(e).Convert(e, node.TypeBool)
					if err != nil {
						unknown = true
						continue
					}
					if item.(*node.Bool).Val == true {
						return item
					}
				}
				if !unknown {
					return &node.Bool{false}
				}
			}
			return node.NewUnaryOp(arg, str)
		},
	)
}

package unary

import (
	"jakobvarmose/compute/node"

	"github.com/ericlagergren/decimal"
	dmath "github.com/ericlagergren/decimal/math"
)

func init() {
	// Trigonometric functions
	CreateUnary(
		"cos",
		func(e node.Environment, arg node.Node, str string) node.Node {
			d, err := arg.Convert(e, node.TypeDecimal)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			val := dmath.Cos(new(decimal.Big), d.(*node.Decimal).Val)
			return &node.Decimal{val}
		},
	)
	CreateUnary(
		"sin",
		func(e node.Environment, arg node.Node, str string) node.Node {
			d, err := arg.Convert(e, node.TypeDecimal)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			val := dmath.Sin(new(decimal.Big), d.(*node.Decimal).Val)
			return &node.Decimal{val}
		},
	)
	CreateUnary(
		"tan",
		func(e node.Environment, arg node.Node, str string) node.Node {
			d, err := arg.Convert(e, node.TypeDecimal)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			val := dmath.Tan(new(decimal.Big), d.(*node.Decimal).Val)
			return &node.Decimal{val}
		},
	)
	CreateUnary(
		"acos",
		func(e node.Environment, arg node.Node, str string) node.Node {
			d, err := arg.Convert(e, node.TypeDecimal)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			val := dmath.Acos(new(decimal.Big), d.(*node.Decimal).Val)
			return &node.Decimal{val}
		},
	)
	CreateUnary(
		"asin",
		func(e node.Environment, arg node.Node, str string) node.Node {
			d, err := arg.Convert(e, node.TypeDecimal)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			val := dmath.Asin(new(decimal.Big), d.(*node.Decimal).Val)
			return &node.Decimal{val}
		},
	)
	CreateUnary(
		"atan",
		func(e node.Environment, arg node.Node, str string) node.Node {
			d, err := arg.Convert(e, node.TypeDecimal)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			val := dmath.Atan(new(decimal.Big), d.(*node.Decimal).Val)
			return &node.Decimal{val}
		},
	)
}

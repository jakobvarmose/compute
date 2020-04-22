package unary

import (
	"jakobvarmose/compute/node"
	"math/big"

	"github.com/ericlagergren/decimal"
	dmath "github.com/ericlagergren/decimal/math"
)

func init() {
	CreateUnary(
		"-",
		func(e node.Environment, arg node.Node, str string) node.Node {
			res, err := arg.Convert(e, node.TypeComplex)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			return res.(*node.Complex).Neg()
		},
	)
	CreateUnary(
		"+",
		func(e node.Environment, arg node.Node, str string) node.Node {
			res, err := arg.Convert(e, node.TypeComplex)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			return res
		},
	)

	// Number functions
	CreateUnary(
		"abs",
		func(e node.Environment, arg node.Node, str string) node.Node {
			switch arg := arg.(type) {
			case *node.Int:
				return &node.Int{new(big.Int).Abs(arg.Big)}
			case *node.Rat:
				return &node.Rat{new(big.Rat).Abs(arg.Big)}
			case *node.Decimal:
				return &node.Decimal{new(decimal.Big).Abs(arg.Val)}
			}
			return node.NewUnaryOp(arg, str)
		},
	)
	CreateUnary(
		"real",
		func(e node.Environment, arg node.Node, str string) node.Node {
			res, err := arg.Convert(e, node.TypeComplex)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			return &node.Rat{res.(*node.Complex).Real}
		},
	)
	CreateUnary(
		"imag",
		func(e node.Environment, arg node.Node, str string) node.Node {
			res, err := arg.Convert(e, node.TypeComplex)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			return &node.Rat{res.(*node.Complex).Imag}
		},
	)
	CreateUnary(
		"conj",
		func(e node.Environment, arg node.Node, str string) node.Node {
			res, err := arg.Convert(e, node.TypeComplex)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			return res.(*node.Complex).Conj()
		},
	)
	CreateUnary(
		"floor",
		func(e node.Environment, arg node.Node, str string) node.Node {
			arg = arg.Eval(e)
			switch arg := arg.(type) {
			case *node.Rat, *node.Int:
				r, _ := arg.Convert(e, node.TypeRat)
				return r.(*node.Rat).Floor()
			}
			return node.NewUnaryOp(arg, str)
		},
	)
	CreateUnary(
		"ceil",
		func(e node.Environment, arg node.Node, str string) node.Node {
			arg = arg.Eval(e)
			switch arg := arg.(type) {
			case *node.Int, *node.Rat:
				r, _ := arg.Convert(e, node.TypeRat)
				return r.(*node.Rat).Ceil()
			}
			return node.NewUnaryOp(arg, str)
		},
	)
	CreateUnary(
		"exp",
		func(e node.Environment, arg node.Node, str string) node.Node {
			d, err := arg.Convert(e, node.TypeDecimal)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			val := dmath.Exp(new(decimal.Big), d.(*node.Decimal).Val)
			return &node.Decimal{val}
		},
	)
	CreateUnary(
		"ln",
		func(e node.Environment, arg node.Node, str string) node.Node {
			d, err := arg.Convert(e, node.TypeDecimal)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			val := dmath.Log(new(decimal.Big), d.(*node.Decimal).Val)
			return &node.Decimal{val}
		},
	)
	CreateUnary(
		"sqrt",
		func(e node.Environment, arg node.Node, str string) node.Node {
			d, err := arg.Convert(e, node.TypeDecimal)
			if err != nil {
				return node.NewUnaryOp(arg, str)
			}
			val := dmath.Sqrt(new(decimal.Big), d.(*node.Decimal).Val)
			return &node.Decimal{val}
		},
	)
}

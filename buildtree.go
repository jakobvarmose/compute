package main

import (
	"errors"

	"jakobvarmose/compute/node"
	"jakobvarmose/compute/token"
)

func buildTree(e node.Environment, tokens <-chan PostfixToken) (node.Node, error) {
	var stack []node.Node
	for t := range tokens {
		switch t.Type {
		case token.TypeInt:
			i, err := node.NewInt(t.Value)
			if err != nil {
				return nil, err
			}
			stack = append(stack, i)
		case token.TypeBool:
			b, err := node.NewBool(t.Value)
			if err != nil {
				return nil, err
			}
			stack = append(stack, b)
		case token.TypeComplex:
			c, err := node.NewComplex(t.Value)
			if err != nil {
				return nil, err
			}
			stack = append(stack, c)
		case token.TypeString:
			s := node.NewString(t.Value)
			stack = append(stack, s)
		case token.TypeVariable:
			stack = append(stack, node.NewVariable(t.Value))
		case token.TypeOperator:
			if len(stack) < t.Arity {
				return nil, errors.New("not enough tokens")
			}
			switch t.Arity {
			case 1:
				a := stack[len(stack)-1]
				stack[len(stack)-1] = node.NewUnaryOp(a, t.Value)
			case 2:
				a := stack[len(stack)-2]
				b := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				switch t.Value {
				case "-":
					stack[len(stack)-1] = node.NewMinus(a, b)
				case "=":
					stack[len(stack)-1] = node.NewAssign(e, a, b)
				case "#":
					stack[len(stack)-1] = node.NewIndex(e, a, b)
				case ":":
					stack[len(stack)-1] = node.NewRegister(e, a, b)
				case "and":
					stack[len(stack)-1] = node.NewAnd(e, a, b)
				case "or":
					stack[len(stack)-1] = node.NewOr(e, a, b)
				default:
					stack[len(stack)-1] = node.NewBinaryOp(a, b, t.Value)
				}
			default:
				return nil, errors.New("wrong arity")
			}
		case token.TypeNewline:
			if len(stack) != 1 {
				return nil, errors.New("syntax error")
			}
			return stack[0], nil
		default:
			return nil, errors.New("unknown token")
		}
	}
	return nil, nil
}

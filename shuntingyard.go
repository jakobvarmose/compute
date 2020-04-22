package main

import (
	"errors"

	"jakobvarmose/compute/node"
	"jakobvarmose/compute/token"
)

type PostfixToken struct {
	Type  token.Type
	Value string
	Arity int
}

func shuntingYard(tokens <-chan token.Token, queue chan<- PostfixToken) error {
	var stack []PostfixToken
	var variable token.Token
	unary := true
	for t := range tokens {
		if variable.Type == token.TypeVariable {
			switch t.Type {
			case token.TypeStart:
				stack = append(stack, PostfixToken{token.TypeOperator, variable.Value, 1})
				variable = token.Token{}
			default:
				queue <- PostfixToken{token.TypeVariable, variable.Value, 0}
				variable = token.Token{}
			case token.TypeSpace:
			}
		}
		switch t.Type {
		case token.TypeInt, token.TypeBool, token.TypeString, token.TypeComplex:
			queue <- PostfixToken{t.Type, t.Value, 0}
		case token.TypeVariable:
			variable = t
		case token.TypeOperator:
			if unary {
				stack = append(stack, PostfixToken{t.Type, t.Value, 1})
			} else {
				for len(stack) > 0 &&
					(stack[len(stack)-1].Arity == 1 ||
						stack[len(stack)-1].Type == token.TypeOperator &&
							(node.LeftAssociative(t.Value) &&
								node.Precedence(t.Value) == node.Precedence(stack[len(stack)-1].Value) ||
								node.Precedence(t.Value) < node.Precedence(stack[len(stack)-1].Value))) {
					queue <- stack[len(stack)-1]
					stack = stack[:len(stack)-1]
				}
				stack = append(stack, PostfixToken{t.Type, t.Value, 2})
			}
		case token.TypeStart:
			stack = append(stack, PostfixToken{t.Type, t.Value, 0})
		case token.TypeEnd:
			// Until we find a left paren
			for len(stack) > 0 &&
				stack[len(stack)-1].Type != token.TypeStart {
				queue <- stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return errors.New("missing left paren")
			}
			// Pop left paren
			stack = stack[:len(stack)-1]
		case token.TypeNewline:
			for len(stack) > 0 &&
				stack[len(stack)-1].Type != token.TypeStart {
				queue <- stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
			if len(stack) != 0 {
				return errors.New("missing right paren")
			}
			queue <- PostfixToken{t.Type, t.Value, 0}
		case token.TypeSpace:
		default:
			panic(t.String())
		}
		unary = t.Type == token.TypeOperator || t.Type == token.TypeStart || t.Type == token.TypeNewline
	}
	return nil
}

package main

import (
	"fmt"
	"io"
	"math/big"
	"os"

	"jakobvarmose/compute/environment"
	"jakobvarmose/compute/lexer"
	"jakobvarmose/compute/node"
	"jakobvarmose/compute/token"
	"jakobvarmose/compute/unary"
)

type Value interface {
	String() string
}

type StringWriter interface {
	Write(e node.Environment, w io.Writer)
}

func parse(in <-chan uint8, out chan<- token.Token) {
	c := <-in
	for {
		if c >= '0' && c <= '9' {
			var str string
			for c >= '0' && c <= '9' {
				str += string(c)
				c = <-in
			}
			out <- token.Token{token.TypeInt, str}
		} else if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
			var str string
			for c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
				str += string(c)
				c = <-in
			}
			out <- token.Token{token.TypeVariable, str}
		} else if c == '+' || c == '-' || c == '*' || c == '/' {
			str := string(c)
			c = <-in
			out <- token.Token{token.TypeOperator, str}
		} else if c == '"' {
			str := string(c)
			c = <-in
			for c != '"' {
				str += string(c)
				c = <-in
			}
			str += string(c)
			c = <-in
			out <- token.Token{token.TypeString, str}
		} else if c == '\n' {
			out <- token.Token{token.TypeNewline, "\n"}
			c = <-in
		} else {
			c = <-in
		}
	}
}

func print(tokens <-chan node.Node) {
	for token := range tokens {
		fmt.Println("Result", token.String())
	}
}

func read(bytes chan<- uint8) {
	b := make([]byte, 1)
	_, err := os.Stdin.Read(b)
	for err == nil {
		bytes <- b[0]
		_, err = os.Stdin.Read(b)
	}
}

func convert(tokens <-chan token.Token, nodes chan<- node.Node) {
	var stack []node.Node
	//var ops []*token.Operator
	for t := range tokens {
		switch t.Type {
		case token.TypeNewline:
			if len(stack) == 1 {
				n := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				nodes <- n
			} else if len(stack) != 0 {
				panic("unexpected newline")
			}
		case token.TypeInt:
			i, err := node.NewInt(t.Value)
			if err != nil {
				panic(err)
			}
			stack = append(stack, i)
		case token.TypeVariable:
			stack = append(stack, node.NewVariable(t.Value))
		case token.TypeString:
			stack = append(stack, node.NewString(t.Value))
		case token.TypeOperator:
			//ops = append(ops, t)
		case token.TypeSpace:
		default:
			panic("unknown token")
		}
	}
}

func Approx(e node.Environment, arg node.Node) (node.Node, error) {
	switch arg.Type() {
	case node.TypeInt, node.TypeRat, node.TypeExact:
		return arg.Convert(e, node.TypeDecimal)
	}
	return arg, nil
}

func main() {
	e := environment.New()
	for _, exact := range node.Exacts {
		e.SetDontSave(exact.Name, exact)
	}
	for _, builtin := range unary.Builtins {
		e.SetDontSave(builtin.Name, builtin)
	}
	/*for _, builtin := range binary.Builtins {
		e.SetDontSave(builtin.Name, builtin)
	}*/

	tokens := lexer.Run(os.Stdin)
	rpnCh := make(chan PostfixToken)
	go func() {
		err := shuntingYard(tokens, rpnCh)
		if err != nil {
			fmt.Println(err)
		}
	}()
	for {
		var root node.Node
		var err error
		var res node.Node
		func() {
			defer func() {
				err := recover()
				if err != nil {
					switch err := err.(type) {
					case error:
						res = node.NewError(err.Error())
					case string:
						res = node.NewError(err)
					default:
						res = node.NewError("unknown recover type")
					}
				}
			}()
			root, err = buildTree(e, rpnCh)
			if err != nil {
				panic(err)
			}
			fmt.Println(root)
			res = root.Eval(e)
		}()
		e.Set("_", res)
		e.Save()
		if b, ok := root.(*node.BinaryOp); ok && b.Name() == "=" {
			continue
		}
		if res.Type() == node.TypeExact || res.Type() == node.TypeRat {
			d, _ := res.Convert(e, node.TypeDecimal)
			fmt.Printf("=> %s ≈ %s\n", res, d)
		} else if i, ok := res.(*node.Int); ok && i.Big.Cmp(big.NewInt(1e7)) >= 0 {
			d, _ := res.Convert(e, node.TypeDecimal)
			fmt.Printf("=> %s ≈ %s\n", res, d)
		} else if _, ok := res.(*node.Builtin); ok {
			fmt.Printf("=> %s <builtin>\n", res)
		} else if w, ok := res.(StringWriter); ok {
			fmt.Print("=> ")
			w.Write(e, os.Stdout)
			fmt.Println()
		} else {
			fmt.Println("=>", res)
		}
	}

	//nodes := make(chan node.Node)
	//go convert(tokens, nodes)
	//print(nodes)
}

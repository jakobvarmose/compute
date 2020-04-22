package unary

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"jakobvarmose/compute/node"
	"math/big"
	"os"
)

func init() {
	CreateUnary(
		"typeof",
		func(e node.Environment, arg node.Node, str string) node.Node {
			return &node.String{arg.ComputedType(e).String()}
		},
	)
	CreateUnary(
		"@",
		func(e node.Environment, arg node.Node, str string) node.Node {
			br := bufio.NewReader(os.Stdin)
			if s, ok := arg.(*node.String); ok {
				fmt.Print(s.Value + ": ")
			}
			line, err := br.ReadString('\n')
			if err != nil {
				return node.NewError(err.Error())
			}
			return &node.String{line[:len(line)-1]}
		},
	)
	CreateUnary(
		"rand",
		func(e node.Environment, arg node.Node, str string) node.Node {
			if true || e.SideEffects() {
				switch arg := arg.(type) {
				case *node.Int:
					if arg.Big.Cmp(big.NewInt(0)) <= 0 {
						return node.NewError(node.NewUnaryOp(arg, str).String() + " is undefined")
					}
					r, err := rand.Int(rand.Reader, arg.Big)
					if err != nil {
						panic(err)
					}
					return &node.Int{r}
				}
			}
			return node.NewUnaryOp(arg, str)
		},
	)
	CreateUnary(
		"name",
		func(e node.Environment, arg node.Node, str string) node.Node {
			if b, ok := arg.(*node.Builtin); ok {
				return &node.String{b.Name}
			}
			return node.NewUnaryOp(arg, str)
		},
	)
	CreateUnary(
		"isPrime",
		func(e node.Environment, arg node.Node, str string) node.Node {
			switch arg := arg.(type) {
			case *node.Int:
				return &node.Bool{isPrime(arg.Big)}
			}
			return node.NewUnaryOp(arg, str)
		},
	)
}

func isPrime(p *big.Int) bool {
	if p.Cmp(big.NewInt(2)) < 0 {
		return false
	}
	if new(big.Int).Mod(p, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		if p.Cmp(big.NewInt(2)) == 0 {
			return true
		}
		return false
	}
	for i := big.NewInt(3); new(big.Int).Mul(i, i).Cmp(p) <= 0; i.Add(i, big.NewInt(2)) {
		if new(big.Int).Mod(p, i).Cmp(big.NewInt(0)) == 0 {
			return false
		}
	}
	return true
}

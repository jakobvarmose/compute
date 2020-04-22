package node

import (
	"math/big"
)

func op_mod(e Environment, b *BinaryOp) Node {
	left := b.left.Eval(e)
	right := b.right.Eval(e)
	switch left := left.(type) {
	case ListInterface:
		return left.Mod(e, right)
	case *Int:
		switch right := b.right.Eval(e).(type) {
		case *Int:
			res := new(big.Int).Rem(left.Big, right.Big)
			return &Int{res}
		}
	}
	return &BinaryOp{left, right, "%"}
}

package node

import (
	"math/big"
)

func ListMod(e Environment, l ListInterface, other Node) Node {
	right := other.Eval(e)
	var list []Node
	buf := e.Get("@")
	for i := big.NewInt(0); i.Cmp(l.Len()) < 0; i = new(big.Int).Add(i, big.NewInt(1)) {
		item := l.ItemAt(i)
		e.Set("@", item)
		list = append(list, right.Eval(e))
	}
	e.Set("@", buf)
	return &List{list, l.ItemType(e)}
}

func ListFilter(e Environment, l ListInterface, other Node) Node {
	right := other.Eval(e)
	var list []Node
	for i := big.NewInt(0); i.Cmp(l.Len()) < 0; i = new(big.Int).Add(i, big.NewInt(1)) {
		item := l.ItemAt(i)
		b, err := (&BinaryOp{right, item, "()"}).Eval(e).Convert(e, TypeBool)
		if err != nil {
			return l
		}
		if b.(*Bool).Val {
			list = append(list, item)
		}
	}
	return &List{list, l.ItemType(e)}
}

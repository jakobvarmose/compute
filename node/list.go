package node

import (
	"errors"
	"io"
	"math/big"

	"jakobvarmose/compute/vtype"
)

type List struct {
	Val      []Node
	itemType vtype.VType
}

func NewList(itemType vtype.VType) *List {
	return &List{nil, itemType}
}

func (l *List) String() string {
	res := "["
	for i, item := range l.Val {
		if i != 0 {
			res += ","
		}
		res += item.String()
	}
	res += "]"
	return res
}
func (l *List) MarshalString() string {
	return "LIST NOT CORRECT"
}
func (l *List) Write(e Environment, w io.Writer) {
	w.Write([]byte("["))
	for i := big.NewInt(0); i.Cmp(l.Len()) < 0; i.Add(i, big.NewInt(1)) {
		if i.Cmp(big.NewInt(0)) != 0 {
			w.Write([]byte(","))
		}
		item := l.ItemAt(i)
		w.Write([]byte(item.String()))
	}
	w.Write([]byte("]"))
}

func (l *List) Eval(e Environment) Node {
	var res []Node
	for _, item := range l.Val {
		res = append(res, item.Eval(e))
	}
	return &List{res, l.itemType}
}

func (l *List) Convert(env Environment, t Type) (Node, error) {
	switch t {
	case TypeList:
		return l, nil
	default:
		return nil, errors.New("cannot convert")
	}
}

func (l *List) Type() Type {
	return TypeList
}

func (l *List) ItemType(e Environment) vtype.VType {
	return l.itemType
}

func (l *List) ComputedType(e Environment) vtype.VType {
	return &vtype.ListType{l.ItemType(e)}
}

func (l *List) Fold(e Environment, res Node, op string) Node {
	for _, item := range l.Val[1:] {
		res = &BinaryOp{res, item.Eval(e), op}
	}
	return res
}

func (l *List) Len() *big.Int {
	return big.NewInt(int64(len(l.Val)))
}

func (l *List) ItemAt(i *big.Int) Node {
	return l.Val[i.Int64()]
}

func (l *List) Split(i *big.Int) (ListInterface, ListInterface) {
	a := &List{l.Val[:i.Int64()], l.itemType}
	b := &List{l.Val[i.Int64():], l.itemType}
	return a, b
}

func (l *List) Mod(e Environment, other Node) Node {
	return ListMod(e, l, other)
}

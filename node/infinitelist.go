package node

import (
	"errors"
	"io"
	"math/big"

	"jakobvarmose/compute/vtype"
)

type ListInterface interface {
	Node

	Len() *big.Int
	ItemAt(i *big.Int) Node
	ItemType(e Environment) vtype.VType
	Split(i *big.Int) (ListInterface, ListInterface)
	Mod(e Environment, other Node) Node
}

type InfiniteList struct {
	size      *big.Int
	statement Node
}

func NewInfiniteList(size *big.Int) *InfiniteList {
	variable := UniqueVariable(vtype.Integer)
	statement := &BinaryOp{
		variable,
		variable,
		"=>",
	}
	return &InfiniteList{size, statement}
}

func (l *InfiniteList) String() string {
	res := "["
	res += l.ItemAt(big.NewInt(0)).String()
	res += ",...,"
	res += l.ItemAt(new(big.Int).Sub(l.size, big.NewInt(1))).String()
	res += "]"
	return res
}
func (l *InfiniteList) MarshalString() string {
	return l.String()
}
func (l *InfiniteList) Write(e Environment, w io.Writer) {
	w.Write([]byte("["))
	if l.size.Cmp(big.NewInt(101)) > 0 {
		for i := big.NewInt(0); i.Cmp(big.NewInt(100)) < 0; i.Add(i, big.NewInt(1)) {
			if i.Cmp(big.NewInt(0)) != 0 {
				w.Write([]byte(","))
			}
			item := l.ItemAt(i).Eval(e)
			w.Write([]byte(item.String()))
		}
		w.Write([]byte(",...,"))
		item := l.ItemAt(new(big.Int).Sub(l.size, big.NewInt(1))).Eval(e)
		w.Write([]byte(item.String()))
		w.Write([]byte("]"))
	} else {
		for i := big.NewInt(0); i.Cmp(l.Len()) < 0; i.Add(i, big.NewInt(1)) {
			if i.Cmp(big.NewInt(0)) != 0 {
				w.Write([]byte(","))
			}
			item := l.ItemAt(i).Eval(e)
			w.Write([]byte(item.String()))
		}
	}
	w.Write([]byte("]"))
}

func (l *InfiniteList) Eval(e Environment) Node {
	return l
}

func (l *InfiniteList) Convert(e Environment, t Type) (Node, error) {
	switch t {
	case TypeList:
		return l, nil
	default:
		return nil, errors.New("cannot convert")
	}
}

func (l *InfiniteList) Type() Type {
	return TypeList
}

func (l *InfiniteList) ItemType(e Environment) vtype.VType {
	return l.statement.ComputedType(e)
}

func (l *InfiniteList) ComputedType(e Environment) vtype.VType {
	return &vtype.ListType{l.ItemType(e)}
}

func (l *InfiniteList) Len() *big.Int {
	return l.size
}

func (l *InfiniteList) ItemAt(i *big.Int) Node {
	return &BinaryOp{
		l.statement,
		&Int{i},
		"()",
	}
}

func (l *InfiniteList) Split(i *big.Int) (ListInterface, ListInterface) {
	a := &InfiniteList{i, l.statement}
	variable := UniqueVariable(vtype.Integer)
	statement := &BinaryOp{
		variable,
		&BinaryOp{
			l.statement,
			&BinaryOp{variable, &Int{big.NewInt(1)}, "+"},
			"()",
		},
		"=>",
	}
	b := &InfiniteList{new(big.Int).Sub(l.size, i), statement}
	return a, b
}

func (l *InfiniteList) Mod(e Environment, other Node) Node {
	variable := UniqueVariable(vtype.Integer)
	statement := &BinaryOp{
		variable,
		&BinaryOp{
			other,
			&BinaryOp{
				l.statement,
				variable,
				"()",
			},
			"()",
		},
		"=>",
	}
	return &InfiniteList{l.size, statement}
}

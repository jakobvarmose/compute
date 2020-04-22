package node

import (
	"errors"
	"jakobvarmose/compute/vtype"
	"math/big"
)

type Tuple struct {
	Items []Node
}

func (t *Tuple) String() string {
	res := "("
	for i, item := range t.Items {
		res += item.String()
		if i < len(t.Items)-1 || len(t.Items) == 1 {
			res += ","
		}
	}
	res += ")"
	return res
}
func (t *Tuple) MarshalString() string {
	return t.String()
}

func (t *Tuple) Eval(e Environment) Node {
	return t
}

func (t *Tuple) Convert(env Environment, typ Type) (Node, error) {
	switch typ {
	case TypeTuple:
		return t, nil
	case TypeString:
		return &String{t.String()}, nil
	default:
		return nil, errors.New("cannot convert")
	}
}

func (t *Tuple) Type() Type {
	return TypeComplex
}

func (t *Tuple) ComputedType(e Environment) vtype.VType {
	return vtype.Boolean
}

func (t *Tuple) Len() *big.Int {
	return big.NewInt(int64(len(t.Items)))
}

func (t *Tuple) ItemAt(i *big.Int) Node {
	return t.Items[int(i.Int64())]
}

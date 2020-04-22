package node

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"jakobvarmose/compute/vtype"
)

type String struct {
	Value string
}

func NewString(str string) *String {
	r := strings.NewReplacer("\\\\", "\\", "\\\"", "\"", "\\n", "\n")
	str = r.Replace(str[1 : len(str)-1])
	return &String{str}
}

func (s *String) String() string {
	return fmt.Sprintf("%q", s.Value)
}
func (s *String) MarshalString() string {
	return s.String()
}

func (s *String) Eval(e Environment) Node {
	return s
}

func (s *String) Convert(env Environment, t Type) (Node, error) {
	switch t {
	case TypeString:
		return s, nil
	default:
		return nil, errors.New("cannot convert")
	}
}

func (s *String) Type() Type {
	return TypeString
}

func (s *String) ComputedType(e Environment) vtype.VType {
	return vtype.String
}

func (s *String) Len() *big.Int {
	return big.NewInt(int64(len(s.Value)))
}

package node

import (
	"errors"
	"jakobvarmose/compute/vtype"
)

type Register struct {
	left  *Variable
	right Node
}

func NewRegister(e Environment, left, right Node) *Register {
	return &Register{left.(*Variable), right}
}

func (op *Register) Name() string {
	return ":"
}
func (op *Register) Left() Node {
	return op.left
}
func (op *Register) Right() Node {
	return op.right
}

func (op *Register) String() string {
	return binaryOpString(op)
}

func (op *Register) MarshalString() string {
	return binaryOpString(op)
}

func (op *Register) Eval(e Environment) Node {
	left := op.left
	if right, ok := op.right.(*Variable); ok {
		switch right.Name {
		case "S":
			e.SetType(left.Name, vtype.String)
		case "B":
			e.SetType(left.Name, vtype.Boolean)

		case "Z":
			e.SetType(left.Name, vtype.Integer)
		case "Q":
			e.SetType(left.Name, vtype.Rational)
		case "R":
			e.SetType(left.Name, vtype.Real)

		case "Zi":
			e.SetType(left.Name, vtype.Gaussian_Integer)
		case "Qi":
			e.SetType(left.Name, vtype.Gaussian_Rational)
		case "C":
			e.SetType(left.Name, vtype.Complex)

		default:
			return &Bool{false}
		}
		return &Bool{true}
	}
	return &Bool{false}
}

func (op *Register) Convert(e Environment, t Type) (Node, error) {
	switch t {
	case TypeUnknown:
		return op, nil
	default:
		return nil, errors.New("cannot convert")
	}
}

func (op *Register) Type() Type {
	return TypeUnknown
}

func (op *Register) ComputedType(e Environment) vtype.VType {
	return vtype.Boolean
}

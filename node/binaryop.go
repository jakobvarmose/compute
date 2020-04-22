package node

import (
	"errors"
	"math/big"
	"strings"

	"jakobvarmose/compute/vtype"
)

type BinaryOp struct {
	left, right Node
	str         string
}

func NewBinaryOp(left, right Node, str string) *BinaryOp {
	return &BinaryOp{left, right, str}
}

func (b *BinaryOp) Name() string {
	return b.str
}
func (b *BinaryOp) Left() Node {
	return b.left
}
func (b *BinaryOp) Right() Node {
	return b.right
}

func binaryOpString(op interface {
	Name() string
	Left() Node
	Right() Node
}) string {
	name := op.Name()
	left := op.Left()
	leftStr := left.String()
	if left, ok := left.(interface{ Name() string }); ok {
		if Precedence(left.Name()) < Precedence(name) ||
			Precedence(left.Name()) == Precedence(name) && !LeftAssociative(name) {
			leftStr = "(" + leftStr + ")"
		}
	}
	right := op.Right()
	rightStr := right.String()
	if right, ok := right.(interface{ Name() string }); ok {
		if Precedence(right.Name()) < Precedence(name) ||
			Precedence(right.Name()) == Precedence(name) && LeftAssociative(name) {
			rightStr = "(" + rightStr + ")"
		}
	}
	if name[0] >= 'a' && name[0] <= 'z' {
		name = " " + name + " "
	}
	return leftStr + name + rightStr
}

func (op *BinaryOp) String() string {
	return binaryOpString(op)
}
func (op *BinaryOp) MarshalString() string {
	return binaryOpString(op)
}

func (b *BinaryOp) Eval(env Environment) Node {
	switch b.str {
	case "%":
		return op_mod(env, b)
	case "=>":
		return op_fn(env, b)
	case "()":
		return op_call(env, b)
	}
	info := operators[b.str]
	if info.f != nil {
		left := b.left.Eval(env)
		right := b.right.Eval(env)
		leftType, rightType := info.whichType(left.Type(), right.Type())
		if leftType == TypeError && rightType == TypeError {
			return &BinaryOp{left, right, b.str}
		}
		left, err := left.Convert(env, leftType)
		if err != nil {
			panic(err)
		}
		right, err = right.Convert(env, rightType)
		if err != nil {
			panic(err)
		}
		f := info.f[leftType]
		if f == nil {
			return &BinaryOp{left, right, b.str}
		}
		return f(env, left, right)
	}
	return NewError("unknown operator")
}

func (b *BinaryOp) Convert(env Environment, t Type) (Node, error) {
	switch t {
	case TypeUnknown:
		return b, nil
	default:
		return nil, errors.New("cannot convert")
	}
}

func (b *BinaryOp) Type() Type {
	return TypeUnknown
}

func (b *BinaryOp) ComputedType(e Environment) vtype.VType {
	switch b.str {
	case "=>":
		return b.right.ComputedType(e)
	}
	return vtype.Boolean
}

type X map[Type]func(e Environment, n1, n2 Node) Node

type operatorInfo struct {
	t         Node
	whichType func(Type, Type) (Type, Type)
	f         X
}

func preferInt(a, b Type) (Type, Type) {
	switch a {
	case TypeInt:
		if b == TypeInt {
			return TypeInt, TypeInt
		} else if b == TypeComplex {
			return TypeComplex, TypeComplex
		} else if b == TypeRat {
			return TypeRat, TypeRat
		}
	case TypeRat:
		if b == TypeComplex {
			return TypeComplex, TypeComplex
		} else if b == TypeInt || b == TypeRat {
			return TypeRat, TypeRat
		}
	case TypeComplex:
		if b == TypeInt || b == TypeRat || b == TypeComplex {
			return TypeComplex, TypeComplex
		}
	}
	return TypeError, TypeError
}

func preferIntAdd(a, b Type) (Type, Type) {
	switch a {
	case TypeInt:
		if b == TypeInt {
			return TypeInt, TypeInt
		} else if b == TypeComplex {
			return TypeComplex, TypeComplex
		} else if b == TypeRat {
			return TypeRat, TypeRat
		}
	case TypeRat:
		if b == TypeComplex {
			return TypeComplex, TypeComplex
		} else if b == TypeInt || b == TypeRat {
			return TypeRat, TypeRat
		}
	case TypeComplex:
		if b == TypeInt || b == TypeRat || b == TypeComplex {
			return TypeComplex, TypeComplex
		}
	case TypeString:
		if b == TypeString {
			return TypeString, TypeString
		}
	}
	return TypeError, TypeError
}

func preferIntMul(a, b Type) (Type, Type) {
	switch a {
	case TypeInt:
		if b == TypeInt {
			return TypeInt, TypeInt
		} else if b == TypeComplex {
			return TypeComplex, TypeComplex
		} else if b == TypeRat {
			return TypeRat, TypeRat
		} else if b == TypeString {
			return TypeInt, TypeString
		}
	case TypeRat:
		if b == TypeComplex {
			return TypeComplex, TypeComplex
		} else if b == TypeInt || b == TypeRat {
			return TypeRat, TypeRat
		}
	case TypeComplex:
		if b == TypeInt || b == TypeRat || b == TypeComplex {
			return TypeComplex, TypeComplex
		}
	case TypeString:
		if b == TypeInt {
			return TypeString, TypeInt
		}
	}
	return TypeError, TypeError
}

func alwaysRat(a, b Type) (Type, Type) {
	switch a {
	case TypeInt, TypeRat:
		if b == TypeInt || b == TypeRat {
			return TypeRat, TypeRat
		} else if b == TypeComplex {
			return TypeComplex, TypeComplex
		}
	case TypeComplex:
		if b == TypeInt || b == TypeRat || b == TypeComplex {
			return TypeComplex, TypeComplex
		}
	case TypeList:
		return TypeList, b
	}
	return TypeError, TypeError
}

func keepType(a, b Type) (Type, Type) {
	if a == TypeBool || b == TypeBool {
		return TypeError, TypeError
	}
	if a == TypeInt {
		a = TypeRat
	}
	if b == TypeInt {
		b = TypeRat
	}
	return a, b
}

func alwaysBool(a, b Type) (Type, Type) {
	if a == TypeBool && b == TypeBool {
		return TypeBool, TypeBool
	}
	return TypeError, TypeError
}

var operators map[string]operatorInfo

func init() {
	operators = map[string]operatorInfo{
		"+": {
			nil,
			preferIntAdd,
			X{
				TypeVariable: func(e Environment, n1, n2 Node) Node {
					v1 := n1.(*Variable)
					v2 := n2.(*Variable)
					if v1.Name == v2.Name {
						return &BinaryOp{&Int{big.NewInt(2)}, v1, "*"}
					}
					if v1.Name < v2.Name {
						return &BinaryOp{v1, v2, "+"}
					}
					return &BinaryOp{v2, v1, "+"}
				},
				TypeInt: func(e Environment, n1, n2 Node) Node {
					i1 := n1.(*Int)
					i2 := n2.(*Int)
					i3 := new(big.Int).Add(i1.Big, i2.Big)
					return &Int{i3}
				},
				TypeRat: func(e Environment, n1, n2 Node) Node {
					r1 := n1.(*Rat)
					r2 := n2.(*Rat)
					r3 := new(big.Rat).Add(r1.Big, r2.Big)
					if r3.Denom().Cmp(big.NewInt(1)) == 0 {
						return &Int{r3.Num()}
					}
					return &Rat{r3}
				},
				TypeComplex: func(e Environment, n1, n2 Node) Node {
					c1 := n1.(*Complex)
					c2 := n2.(*Complex)
					return c1.Add(c2)
				},
				TypeString: func(e Environment, n1, n2 Node) Node {
					s1 := n1.(*String)
					s2 := n2.(*String)
					return &String{s1.Value + s2.Value}
				},
			},
		},
		"-": {
			nil,
			preferInt,
			X{
				TypeInt: func(e Environment, n1, n2 Node) Node {
					i1 := n1.(*Int)
					i2 := n2.(*Int)
					i3 := new(big.Int).Sub(i1.Big, i2.Big)
					return &Int{i3}
				},
				TypeRat: func(e Environment, n1, n2 Node) Node {
					r1 := n1.(*Rat)
					r2 := n2.(*Rat)
					r3 := new(big.Rat).Sub(r1.Big, r2.Big)
					if r3.Denom().Cmp(big.NewInt(1)) == 0 {
						return &Int{r3.Num()}
					}
					return &Rat{r3}
				},
				TypeComplex: func(e Environment, n1, n2 Node) Node {
					c1 := n1.(*Complex)
					c2 := n2.(*Complex)
					return c1.Sub(c2)
				},
			},
		},
		"*": {
			nil,
			preferIntMul,
			X{
				TypeInt: func(e Environment, n1, n2 Node) Node {
					i1 := n1.(*Int)
					switch i2 := n2.(type) {
					case *Int:
						i3 := new(big.Int).Mul(i1.Big, i2.Big)
						return &Int{i3}
					case *String:
						if !i1.Big.IsUint64() && i1.Big.Uint64() < 1000000 {
							return NewError("string repeat too large")
						}
						return &String{strings.Repeat(i2.Value, int(i1.Big.Uint64()))}
					default:
						panic("wrong type")
					}
				},
				TypeRat: func(e Environment, n1, n2 Node) Node {
					r1 := n1.(*Rat)
					r2 := n2.(*Rat)
					r3 := new(big.Rat).Mul(r1.Big, r2.Big)
					if r3.Denom().Cmp(big.NewInt(1)) == 0 {
						return &Int{r3.Num()}
					}
					return &Rat{r3}
				},
				TypeComplex: func(e Environment, n1, n2 Node) Node {
					c1 := n1.(*Complex)
					c2 := n2.(*Complex)
					return c1.Mul(c2)
				},
				TypeString: func(e Environment, n1, n2 Node) Node {
					s1 := n1.(*String)
					i2 := n2.(*Int)
					if !i2.Big.IsUint64() && i2.Big.Int64() < 1000000 {
						return NewError("string repeat too large")
					}
					return &String{strings.Repeat(s1.Value, int(i2.Big.Uint64()))}
				},
			},
		},
		"/": {
			nil,
			alwaysRat,
			X{
				TypeRat: func(e Environment, n1, n2 Node) Node {
					r1 := n1.(*Rat)
					r2 := n2.(*Rat)
					if r2.Big.Cmp(big.NewRat(0, 1)) == 0 {
						return NewError("division by zero")
					}
					r3 := new(big.Rat).Quo(r1.Big, r2.Big)
					if r3.Denom().Cmp(big.NewInt(1)) == 0 {
						return &Int{r3.Num()}
					}
					return &Rat{r3}
				},
				TypeComplex: func(e Environment, n1, n2 Node) Node {
					c1 := n1.(*Complex)
					c2 := n2.(*Complex)
					return c1.Div(c2)
				},
				TypeList: func(e Environment, n1, n2 Node) Node {
					return ListFilter(e, n1.(ListInterface), n2).Eval(e)
				},
			},
		},
		"**": {
			nil,
			preferInt,
			X{
				TypeInt: func(e Environment, n1, n2 Node) Node {
					i1 := n1.(*Int)
					i2 := n2.(*Int)
					// FIXME handle i2 < 0
					if i2.Big.Cmp(big.NewInt(0)) < 0 {
						panic("not implemented")
					}
					i3 := new(big.Int).Exp(i1.Big, i2.Big, nil)
					return &Int{i3}
				},
			},
		},
		"==": {
			nil,
			keepType,
			X{
				TypeBuiltin: func(e Environment, n1, n2 Node) Node {
					s1 := n1.(*Builtin)
					s2 := n2.(*Builtin)
					return &Bool{s1.Name == s2.Name}
				},
				TypeString: func(e Environment, n1, n2 Node) Node {
					s1 := n1.(*String)
					s2 := n2.(*String)
					return &Bool{s1.Value == s2.Value}
				},
				TypeRat: func(e Environment, n1, n2 Node) Node {
					r1 := n1.(*Rat)
					r2 := n2.(*Rat)
					return &Bool{r1.Big.Cmp(r2.Big) == 0}
				},
			},
		},
		"!=": {
			nil,
			keepType,
			X{
				TypeString: func(e Environment, n1, n2 Node) Node {
					s1 := n1.(*String)
					s2 := n2.(*String)
					return &Bool{s1.Value != s2.Value}
				},
				TypeRat: func(e Environment, n1, n2 Node) Node {
					r1 := n1.(*Rat)
					r2 := n2.(*Rat)
					return &Bool{r1.Big.Cmp(r2.Big) != 0}
				},
			},
		},
		"<": {
			nil,
			keepType,
			X{
				TypeRat: func(e Environment, n1, n2 Node) Node {
					r1 := n1.(*Rat)
					r2 := n2.(*Rat)
					return &Bool{r1.Big.Cmp(r2.Big) < 0}
				},
			},
		},
		"<=": {
			nil,
			keepType,
			X{
				TypeRat: func(e Environment, n1, n2 Node) Node {
					r1 := n1.(*Rat)
					r2 := n2.(*Rat)
					return &Bool{r1.Big.Cmp(r2.Big) <= 0}
				},
			},
		},
		">": {
			nil,
			keepType,
			X{
				TypeRat: func(e Environment, n1, n2 Node) Node {
					r1 := n1.(*Rat)
					r2 := n2.(*Rat)
					return &Bool{r1.Big.Cmp(r2.Big) > 0}
				},
			},
		},
		">=": {
			nil,
			keepType,
			X{
				TypeRat: func(e Environment, n1, n2 Node) Node {
					r1 := n1.(*Rat)
					r2 := n2.(*Rat)
					return &Bool{r1.Big.Cmp(r2.Big) >= 0}
				},
			},
		},
		"and": {
			nil,
			alwaysBool,
			X{
				TypeBool: func(e Environment, n1, n2 Node) Node {
					b1 := n1.(*Bool)
					b2 := n2.(*Bool)
					return &Bool{b1.Val && b2.Val}
				},
			},
		},
		"or": {
			nil,
			alwaysBool,
			X{
				TypeBool: func(e Environment, n1, n2 Node) Node {
					b1 := n1.(*Bool)
					b2 := n2.(*Bool)
					return &Bool{b1.Val || b2.Val}
				},
			},
		},
	}
}

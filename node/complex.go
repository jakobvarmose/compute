package node

import (
	"errors"
	"math/big"
	"strings"

	"jakobvarmose/compute/vtype"
)

type Complex struct {
	Real, Imag *big.Rat
}

func NewComplex(str string) (*Complex, error) {
	if str[len(str)-1] == 'i' {
		r, ok := new(big.Rat).SetString(str[:len(str)-1])
		if !ok {
			return nil, errors.New("invalid complex")
		}
		return &Complex{new(big.Rat), r}, nil
	}
	r, ok := new(big.Rat).SetString(str)
	if !ok {
		return nil, errors.New("invalid complex")
	}
	return &Complex{r, new(big.Rat)}, nil
}

func (c *Complex) String() string {
	real := c.Real.String()
	imag := c.Imag.String()
	if strings.HasSuffix(real, "/1") {
		real = real[:len(real)-2]
	}
	if strings.HasSuffix(imag, "/1") {
		imag = imag[:len(imag)-2]
	} else {
		imag = "(" + imag + ")"
	}
	if imag == "0" {
		return real
	}
	if real == "0" {
		return imag + "i"
	}
	if !strings.HasPrefix(imag, "-") {
		imag = "+" + imag
	}
	return "(" + real + imag + "i)"
}
func (c *Complex) MarshalString() string {
	return "(" + c.Real.RatString() + "+(" + c.Imag.RatString() + ")*1i)"
}

func (c *Complex) Eval(e Environment) Node {
	return c
}

func (c *Complex) Convert(env Environment, t Type) (Node, error) {
	switch t {
	case TypeComplex:
		return c, nil
	case TypeString:
		return &String{c.String()}, nil
	default:
		return nil, errors.New("cannot convert")
	}
}

func (c *Complex) Type() Type {
	return TypeComplex
}

func simplestComplex(real, imag *big.Rat) Node {
	if imag.Cmp(new(big.Rat)) != 0 {
		return &Complex{real, imag}
	}
	if real.Denom().Cmp(big.NewInt(1)) != 0 {
		return &Rat{real}
	}
	return &Int{real.Num()}
}

func (c1 *Complex) Add(c2 *Complex) Node {
	real := new(big.Rat).Add(c1.Real, c2.Real)
	imag := new(big.Rat).Add(c1.Imag, c2.Imag)
	return simplestComplex(real, imag)
}

func (c1 *Complex) Sub(c2 *Complex) Node {
	real := new(big.Rat).Sub(c1.Real, c2.Real)
	imag := new(big.Rat).Sub(c1.Imag, c2.Imag)
	return simplestComplex(real, imag)
}

func (c1 *Complex) Mul(c2 *Complex) Node {
	real := new(big.Rat).Mul(c1.Real, c2.Real)
	real.Sub(real, new(big.Rat).Mul(c1.Imag, c2.Imag))

	imag := new(big.Rat).Mul(c1.Real, c2.Imag)
	imag.Add(imag, new(big.Rat).Mul(c1.Imag, c2.Real))

	return simplestComplex(real, imag)
}

func (c1 *Complex) Div(c2 *Complex) Node {
	denom := new(big.Rat).Mul(c2.Real, c2.Real)
	denom.Add(denom, new(big.Rat).Mul(c2.Imag, c2.Imag))
	denom.Inv(denom)

	real := new(big.Rat).Mul(c1.Real, c2.Real)
	real.Add(real, new(big.Rat).Mul(c1.Imag, c2.Imag))
	real.Mul(real, denom)

	imag := new(big.Rat).Mul(c1.Imag, c2.Real)
	imag.Sub(imag, new(big.Rat).Mul(c1.Real, c2.Imag))
	imag.Mul(imag, denom)

	return simplestComplex(real, imag)
}

func (c *Complex) Neg() Node {
	real := new(big.Rat).Neg(c.Real)
	imag := new(big.Rat).Neg(c.Imag)

	return simplestComplex(real, imag)
}

func (c *Complex) RealRat() Node {
	real := c.Real
	imag := new(big.Rat)

	return simplestComplex(real, imag)
}
func (c *Complex) ImagRat() Node {
	real := c.Imag
	imag := new(big.Rat)

	return simplestComplex(real, imag)
}
func (c *Complex) Conj() Node {
	real := c.Real
	imag := new(big.Rat).Neg(c.Imag)

	return simplestComplex(real, imag)
}

func (c *Complex) ComputedType(e Environment) vtype.VType {
	return vtype.Complex
}

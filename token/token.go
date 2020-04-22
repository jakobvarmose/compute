package token

import (
	"fmt"
)

type Type int

const (
	TypeEOF Type = iota
	TypeItem
	TypeOperator
	TypeVariable
	TypeInt
	TypeComplex
	TypeNewline
	TypeString
	TypeSpace
	TypeStart
	TypeEnd
	TypeBool
	TypeFunction
	TypeLeftSquare
	TypeRightSquare
)

var typeString = map[Type]string{
	TypeEOF:      "EOF",
	TypeInt:      "Int",
	TypeVariable: "Variable",
	TypeNewline:  "Newline",
	TypeSpace:    "Space",
	TypeOperator: "Operator",
	TypeStart:    "Start",
	TypeEnd:      "End",
	TypeBool:     "Bool",
}

type Token struct {
	Type  Type
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("[%s %s]", typeString[t.Type], t.Value)
}

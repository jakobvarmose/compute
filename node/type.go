package node

type Type int

const (
	TypeBool Type = iota
	TypeError
	TypeInt
	TypeNode
	TypeRat
	TypeString
	TypeVariable
	TypeUnknown
	TypeComplex
	TypeList
	TypeDecimal
	TypeExact
	TypeBuiltin
	TypeTuple
)

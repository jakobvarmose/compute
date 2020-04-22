package vtype

type VType interface {
	String() string
}

var (
	String  = &SimpleType{"S"}
	Boolean = &SimpleType{"B"}

	Integer  = &SimpleType{"Z"}
	Rational = &SimpleType{"Q"}
	Real     = &SimpleType{"R"}

	Gaussian_Integer  = &SimpleType{"Zi"}
	Gaussian_Rational = &SimpleType{"Qi"}
	Complex           = &SimpleType{"C"}
	Error             = &SimpleType{"E"}
)

type SimpleType struct {
	Name string
}

func (st *SimpleType) String() string {
	return st.Name
}

type TupleType struct {
	Items []VType
}

func (tt *TupleType) String() string {
	res := "("
	for i, item := range tt.Items {
		if i > 0 {
			res += "*"
		}
		res += item.String()
	}
	res += ")"
	return res
}

type ListType struct {
	Item VType
}

func (lt *ListType) String() string {
	return lt.Item.String() + "[]"
}

type ConvertType struct {
	Left  VType
	Right VType
}

func (ct *ConvertType) String() string {
	return ct.Left.String() + "->" + ct.Right.String()
}

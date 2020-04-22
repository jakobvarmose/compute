package node

func precedence(op string) (int, bool) {
	switch op {
	case "()":
		return 100, true
	case "#":
		return 70, true
	case "**":
		return 60, false
	case "*", "/", "%":
		return 50, true
	case "+", "-":
		return 40, true
	case "==", "!=", "<", "<=", ">", ">=":
		return 30, true
	case "and":
		return 25, true
	case "or":
		return 20, true
	case "=>":
		return 15, false
	case "=":
		return 10, false
	}
	panic("unknown precedence or associativity \"" + op + "\"")
}

func Precedence(op string) int {
	p, _ := precedence(op)
	return p
}

func LeftAssociative(op string) bool {
	_, a := precedence(op)
	return a
}

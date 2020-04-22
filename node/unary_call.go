package node

func unary_call(e Environment, u *UnaryOp) Node {
	fn := e.Get(u.str)
	if fn != nil {
		return (&BinaryOp{fn, u.arg, "()"}).Eval(e)
	}
	return u
}

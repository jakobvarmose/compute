package node

func op_call(e Environment, b *BinaryOp) Node {
	if fn, ok := b.left.(*BinaryOp); ok {
		if fn.str == "=>" {
			arg := b.right
			if left, ok := fn.left.(*Variable); ok {
				buf := e.Get(left.Name)
				e.Set(left.Name, arg)
				res := fn.right.Eval(e)
				e.Set(left.Name, buf)
				return res
			}
		}
	}
	if fn, ok := b.left.(*Builtin); ok {
		return fn.Fn(e, b.right.Eval(e), fn.Name)
	}
	return &BinaryOp{b.left, b.right.Eval(e), b.str}
}

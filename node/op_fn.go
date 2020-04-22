package node

func op_fn(e Environment, b *BinaryOp) Node {
	switch left := b.left.(type) {
	case *Variable:
		buf := e.Get(left.Name)
		e.Set(left.Name, nil)
		res := b.right.Eval(e)
		e.Set(left.Name, buf)
		return &BinaryOp{left, res, "=>"}
	}
	return b
}

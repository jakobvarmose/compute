package node

func op_div(e Environment, b *BinaryOp) Node {
	left := b.left.Eval(e)

	switch left := left.(type) {
	case ListInterface:
		return ListFilter(e, left, b.right).Eval(e)
	}
	right := b.right.Eval(e)
	return &BinaryOp{left, right, b.str}
}

package node

func op_and(e Environment, b *BinaryOp) Node {
	ch := make(chan Node)
	go func() {
		ch <- b.left.Eval(e)
	}()
	go func() {
		ch <- b.right.Eval(e)
	}()
	first := <-ch
	if bo, ok := first.(*Bool); ok {
		if bo.Val == false {
			return &Bool{false}
		}
		second := <-ch
		if bo, ok := second.(*Bool); ok {
			return &Bool{bo.Val}
		}
		return &BinaryOp{first, second, b.str}
	} else {
		second := <-ch
		if bo, ok := second.(*Bool); ok {
			if bo.Val == false {
				return &Bool{false}
			}
		}
		return &BinaryOp{first, second, b.str}
	}
}

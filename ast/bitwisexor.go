package ast

import (
	"fmt"

	"github.com/NicoNex/tau/obj"
)

type BitwiseXor struct {
	l Node
	r Node
}

func NewBitwiseXor(l, r Node) Node {
	return BitwiseXor{l, r}
}

func (b BitwiseXor) Eval(env *obj.Env) obj.Object {
	var left = b.l.Eval(env)
	var right = b.r.Eval(env)

	if isError(left) {
		return left
	}
	if isError(right) {
		return right
	}

	if !assertTypes(left, obj.INT) {
		return obj.NewError("unsupported operator '^' for type %v", left.Type())
	}
	if !assertTypes(right, obj.INT) {
		return obj.NewError("unsupported operator '^' for type %v", right.Type())
	}
	l := left.(*obj.Integer).Val()
	r := right.(*obj.Integer).Val()
	return obj.NewInteger(l ^ r)
}

func (b BitwiseXor) String() string {
	return fmt.Sprintf("(%v ^ %v)", b.l, b.r)
}

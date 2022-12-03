package ast

import (
	"fmt"

	"github.com/NicoNex/tau/internal/code"
	"github.com/NicoNex/tau/internal/compiler"
	"github.com/NicoNex/tau/internal/obj"
)

type BitwiseNot struct {
	n   Node
	pos int
}

func NewBitwiseNot(n Node, pos int) Node {
	return BitwiseNot{
		n:   n,
		pos: pos,
	}
}

func (b BitwiseNot) Eval(env *obj.Env) obj.Object {
	var value = obj.Unwrap(b.n.Eval(env))

	if takesPrecedence(value) {
		return value
	}

	if !obj.AssertTypes(value, obj.IntType) {
		return obj.NewError("unsupported operator '~' for type %v", value.Type())
	}

	n := value.(obj.Integer)
	return obj.Integer(^n)
}

func (b BitwiseNot) String() string {
	return fmt.Sprintf("~%v", b.n)
}

func (b BitwiseNot) Compile(c *compiler.Compiler) (position int, err error) {
	if b.IsConstExpression() {
		o := b.Eval(nil)
		if e, ok := o.(obj.Error); ok {
			return 0, c.NewError(b.pos, string(e))
		}
		position = c.Emit(code.OpConstant, c.AddConstant(o))
		c.Bookmark(b.pos)
		return
	}

	if position, err = b.n.Compile(c); err != nil {
		return
	}
	position = c.Emit(code.OpBwNot)
	c.Bookmark(b.pos)
	return
}

func (b BitwiseNot) IsConstExpression() bool {
	return b.n.IsConstExpression()
}

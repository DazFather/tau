package ast

import (
	"fmt"

	"github.com/NicoNex/tau/internal/compiler"
	"github.com/NicoNex/tau/internal/obj"
)

type MinusMinus struct {
	r   Node
	pos int
}

func NewMinusMinus(r Node, pos int) Node {
	return MinusMinus{
		r:   r,
		pos: pos,
	}
}

func (m MinusMinus) Eval(env *obj.Env) obj.Object {
	var (
		name  string
		right = m.r.Eval(env)
	)

	if takesPrecedence(right) {
		return right
	}

	if ident, ok := m.r.(Identifier); ok {
		name = ident.String()
	}

	if !obj.AssertTypes(right, obj.IntType, obj.FloatType) {
		return obj.NewError("unsupported operator '--' for type %v", right.Type())
	}

	if obj.AssertTypes(right, obj.IntType) {
		if gs, ok := right.(obj.GetSetter); ok {
			r := gs.Object().(obj.Integer)
			return gs.Set(obj.Integer(r - 1))
		}

		r := right.(obj.Integer)
		return env.Set(name, obj.Integer(r-1))
	}

	if gs, ok := right.(obj.GetSetter); ok {
		rightFl, _ := obj.ToFloat(gs.Object(), obj.NullObj)
		r := rightFl.(obj.Float)
		return gs.Set(obj.Float(r - 1))
	}

	rightFl, _ := obj.ToFloat(right, obj.NullObj)
	r := rightFl.(obj.Float)
	return env.Set(name, obj.Float(r-1))
}

func (m MinusMinus) String() string {
	return fmt.Sprintf("--%v", m.r)
}

func (m MinusMinus) Compile(c *compiler.Compiler) (position int, err error) {
	n := Assign{l: m.r, r: Minus{l: m.r, r: Integer(1), pos: m.pos}, pos: m.pos}
	position, err = n.Compile(c)
	c.Bookmark(n.pos)
	return
}

func (m MinusMinus) IsConstExpression() bool {
	return false
}

package ast

import (
	"fmt"

	"github.com/NicoNex/tau/code"
	"github.com/NicoNex/tau/compiler"
	"github.com/NicoNex/tau/obj"
)

type Assign struct {
	l Node
	r Node
}

type setter interface {
	Set(obj.Object) obj.Object
}

func NewAssign(l, r Node) Node {
	return Assign{l, r}
}

func (a Assign) Eval(env *obj.Env) obj.Object {
	if left, ok := a.l.(Identifier); ok {
		right := a.r.Eval(env)
		if takesPrecedence(right) {
			return right
		}
		return env.Set(left.String(), right)
	}

	left := a.l.Eval(env)
	if s, ok := left.(setter); ok {
		right := a.r.Eval(env)
		if takesPrecedence(right) {
			return right
		}
		return s.Set(right)
	}

	return obj.NewError("cannot assign to literal")
}

func (a Assign) String() string {
	return fmt.Sprintf("(%v = %v)", a.l, a.r)
}

func (a Assign) Compile(c *compiler.Compiler) (position int, err error) {
	if left, ok := a.l.(Identifier); ok {
		if position, err = a.r.Compile(c); err != nil {
			return
		}

		symbol := c.Define(string(left))
		if symbol.Scope == compiler.GlobalScope {
			return c.Emit(code.OpSetGlobal, symbol.Index), nil
		} else {
			return c.Emit(code.OpSetLocal, symbol.Index), nil
		}
	}
	return 0, fmt.Errorf("cannot assign to literal")
}

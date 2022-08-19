package ast

import (
	"fmt"

	"github.com/NicoNex/tau/internal/code"
	"github.com/NicoNex/tau/internal/compiler"
	"github.com/NicoNex/tau/internal/obj"
)

type RawString string

func NewRawString(s string) Node {
	return RawString(s)
}

func (r RawString) Eval(env *obj.Env) obj.Object {
	return obj.NewString(string(r))
}

func (r RawString) String() string {
	return string(r)
}

func (r RawString) Quoted() string {
	return fmt.Sprintf("`%s`", string(r))
}

func (r RawString) Compile(c *compiler.Compiler) (position int, err error) {
	return c.Emit(code.OpConstant, c.AddConstant(obj.NewString(string(r)))), nil
}

func (r RawString) Format(prefix string) string {
	return fmt.Sprintf("%s`%s`", prefix, string(r))
}

package ast

import (
	"fmt"
	"sort"
	"strings"

	"github.com/NicoNex/tau/internal/code"
	"github.com/NicoNex/tau/internal/compiler"
	"github.com/NicoNex/tau/internal/obj"
)

type Map struct {
	m   [][2]Node
	pos int
}

func NewMap(pos int, pairs ...[2]Node) Node {
	return Map{
		m:   pairs,
		pos: pos,
	}
}

func (m Map) Eval(env *obj.Env) obj.Object {
	var ret = obj.NewMap()

	for _, pair := range m.m {
		var key, val = pair[0], pair[1]

		k := obj.Unwrap(key.Eval(env))
		if takesPrecedence(k) {
			return k
		}

		h, ok := k.(obj.Hashable)
		if !ok {
			return obj.NewError("invalid map key type %v", k.Type())
		}

		v := obj.Unwrap(val.Eval(env))
		if takesPrecedence(v) {
			return v
		}

		ret.Set(h.KeyHash(), obj.MapPair{Key: k, Value: v})
	}

	return ret
}

func (m Map) String() string {
	var (
		buf strings.Builder
		i   = 1
	)

	buf.WriteString("{")
	for _, pair := range m.m {
		var (
			k   = pair[0]
			v   = pair[1]
			key string
			val string
		)

		if s, ok := k.(String); ok {
			key = s.Quoted()
		} else {
			key = k.String()
		}

		if s, ok := v.(String); ok {
			val = s.Quoted()
		} else {
			val = v.String()
		}

		buf.WriteString(fmt.Sprintf("%s: %s", key, val))

		if i < len(m.m) {
			buf.WriteString(", ")
		}
		i += 1
	}
	buf.WriteString("}")
	return buf.String()
}

func (m Map) Compile(c *compiler.Compiler) (position int, err error) {
	sort.Slice(m.m, func(i, j int) bool {
		return m.m[i][0].String() < m.m[j][0].String()
	})

	for _, pair := range m.m {
		if position, err = pair[0].Compile(c); err != nil {
			return
		}

		if position, err = pair[1].Compile(c); err != nil {
			return
		}
	}

	position = c.Emit(code.OpMap, len(m.m)*2)
	c.Bookmark(m.pos)
	return
}

func (m Map) IsConstExpression() bool {
	return false
}

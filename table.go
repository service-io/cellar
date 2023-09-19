// Package cellar
// @author tabuyos
// @since 2023/9/11
// @description cellar
package cellar

import (
	"github.com/service-io/cellar/token"
	"strings"
)

type RefTable struct {
	refs         []*RefTable
	name         string
	nameTask     Task
	aliasTask    Task
	decorateTask Task
	onTask       Task
	refTask      Task
	typeTask     Task
	renderTask   Task
}

func WithTable(name string) *RefTable {
	return &RefTable{
		name: name,
		nameTask: func(builders ...*strings.Builder) []any {
			for _, builder := range builders {
				builder.WriteString(token.Space.Literal())
				builder.WriteString(name)
			}
			return nil
		},
		decorateTask: func(builders ...*strings.Builder) []any {
			for _, builder := range builders {
				builder.WriteString(name)
			}
			return nil
		},
	}
}

func (t *RefTable) Self() string {
	return t.name
}

func (t *RefTable) As(as string) *RefTable {
	t.decorateTask = func(builders ...*strings.Builder) []any {
		for _, builder := range builders {
			builder.WriteString(as)
		}
		return nil
	}
	t.aliasTask = func(builders ...*strings.Builder) []any {
		for _, builder := range builders {
			builder.WriteString(token.As.Pretty())
			builder.WriteString(as)
		}
		return nil
	}
	return t
}

func (t *RefTable) On(at Selfish, sm Operator, bt Selfish) *RefTable {
	t.onTask = func(buffers ...*strings.Builder) []any {
		buf := buffers[0]
		snip := strings.Join([]string{token.On.Literal(), at.Self(), sm.Self(), bt.Self()}, token.Space.Literal())
		buf.WriteString(token.Space.Literal())
		buf.WriteString(snip)
		return nil
	}
	return t
}

func (t *RefTable) OnEQ(at Selfish, bt Selfish) *RefTable {
	return t.On(at, EQOpr, bt)
}

func (t *RefTable) Ref(refs ...*RefTable) *RefTable {
	t.refs = refs
	t.refTask = func(builders ...*strings.Builder) []any {
		snipBuilder, release := getBuilder()
		defer release()
		for _, ref := range refs {
			snipBuilder.WriteString(token.Space.Literal())
			ref.Render(snipBuilder)
		}

		for _, builder := range builders {
			builder.WriteString(snipBuilder.String())
		}
		return nil
	}
	return t
}

func (t *RefTable) Decorate(key string) string {
	if t.decorateTask == nil {
		return key
	}
	builder, release := getBuilder()
	defer release()
	t.decorateTask.Idle(builder)
	return builder.String() + token.Dot.Literal() + key
}

func (t *RefTable) JoinType(jt JoinType) *RefTable {
	t.typeTask = func(builders ...*strings.Builder) []any {
		snip := jt.Literal()
		for _, buffer := range builders {
			buffer.WriteString(snip)
		}
		return nil
	}
	return t
}

func (t *RefTable) Join() *RefTable {
	return t.JoinType(Join)
}

func (t *RefTable) LeftJoin() *RefTable {
	return t.JoinType(LeftJoin)
}

func (t *RefTable) RightJoin() *RefTable {
	return t.JoinType(RightJoin)
}

func (t *RefTable) FlatAll() (tables []*RefTable) {
	if t == nil {
		return nil
	}
	tables = append(tables, t)
	for _, ref := range t.refs {
		flatAll := ref.FlatAll()
		tables = append(tables, flatAll...)
	}
	return
}

func (t *RefTable) Literal() string {
	builder, release := getBuilder()
	defer release()
	t.Render(builder)
	return builder.String()
}

func (t *RefTable) Render(builders ...*strings.Builder) []any {
	var values []any
	values = append(values, t.typeTask.Idle(builders...)...)
	values = append(values, t.nameTask.Idle(builders...)...)
	values = append(values, t.aliasTask.Idle(builders...)...)
	values = append(values, t.refTask.Idle(builders...)...)
	values = append(values, t.onTask.Idle(builders...)...)
	return values
}

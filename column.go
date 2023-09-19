// Package cellar
// @author tabuyos
// @since 2023/9/11
// @description cellar
package cellar

import (
	"fmt"
	"github.com/service-io/cellar/token"
)

type Column[T any] struct {
	name         string
	aliasFunc    func() string
	mapperFunc   func(*T) any
	decorateFunc func(string) string
}

func WithColumn[T any](name string, mapperFunc func(*T) any, decorateFunc func(string) string) *Column[T] {
	return &Column[T]{name: name, mapperFunc: mapperFunc, decorateFunc: decorateFunc}
}

func (c *Column[T]) Name() string {
	if c.decorateFunc == nil {
		return c.name
	}
	return c.decorateFunc(c.name)
}

func (c *Column[T]) Mapper() func(*T) any {
	return c.mapperFunc
}

func (c *Column[T]) As(as string) *Column[T] {
	c.aliasFunc = func() string {
		return as
	}
	return c
}

func (c *Column[T]) Alias() string {
	if c.aliasFunc == nil {
		return ""
	}
	return c.aliasFunc()
}

func (c *Column[T]) Decorate(fn func(string) string) *Column[T] {
	c.decorateFunc = fn
	return c
}

func (c *Column[T]) Apply(fn func(string) string) *Column[T] {
	c.name = fn(c.name)
	return c
}

func (c *Column[T]) Literal() string {
	if c.aliasFunc == nil {
		return c.Name()
	}
	return c.Name() + token.As.Pretty() + c.aliasFunc()
}

func (c *Column[T]) Self() string {
	return c.Name()
}

func render[T any](asc bool, col *Column[T]) *Order {
	return &Order{
		col: col.Name(),
		asc: asc,
	}
}

func (c *Column[T]) Desc() *Order {
	return render(false, c)
}

func (c *Column[T]) Asc() *Order {
	return render(true, c)
}

// ============================== PRED ==============================

// Cond operator(or fun): cond
func (c *Column[T]) Cond(sym Operator, val ...any) *Predicate {
	return WithPred(c, sym, val...)
}

// IN operator(or fun): in
func (c *Column[T]) IN(values ...any) *Predicate {
	return c.Cond(GenINOpr(len(values)), values...)
}

// EQ operator(or fun): eq
func (c *Column[T]) EQ(values ...any) *Predicate {
	return c.Cond(EQOpr, values...)
}

// NQ operator(or fun): nq
func (c *Column[T]) NQ(values ...any) *Predicate {
	return c.Cond(NQOpr, values...)
}

// LT operator(or fun): lt
func (c *Column[T]) LT(values ...any) *Predicate {
	return c.Cond(LTOpr, values...)
}

// GT operator(or fun): gt
func (c *Column[T]) GT(values ...any) *Predicate {
	return c.Cond(GTOpr, values...)
}

// GE operator(or fun): ge
func (c *Column[T]) GE(values ...any) *Predicate {
	return c.Cond(GEOpr, values...)
}

// LE operator(or fun): le
func (c *Column[T]) LE(values ...any) *Predicate {
	return c.Cond(LEOpr, values...)
}

// Between operator(or fun): between
func (c *Column[T]) Between(values ...any) *Predicate {
	return c.Cond(BetweenOpr, values...)
}

// NotBetween operator(or fun): between
func (c *Column[T]) NotBetween(values ...any) *Predicate {
	return c.Cond(NotBetweenOpr, values...)
}

// Multi operator(or fun): multi
func (c *Column[T]) Multi(values ...any) *Predicate {
	return c.Cond(MultiOpr, values...)
}

// Add operator(or fun): add
func (c *Column[T]) Add(values ...any) *Predicate {
	return c.Cond(AddOpr, values...)
}

// Minus operator(or fun): minus
func (c *Column[T]) Minus(values ...any) *Predicate {
	return c.Cond(MinusOpr, values...)
}

// Like operator(or fun): like
func (c *Column[T]) Like(values ...any) *Predicate {
	return c.Cond(LikeOpr, values...)
}

// LikeLeft operator(or fun): like %xx
func (c *Column[T]) LikeLeft(values ...any) *Predicate {
	newValues := make([]any, len(values))
	for i, a := range values {
		newValues[i] = fmt.Sprintf("%%%v", a)
	}
	return c.Like(newValues...)
}

// LikeRight operator(or fun): like xx%
func (c *Column[T]) LikeRight(values ...any) *Predicate {
	newValues := make([]any, len(values))
	for i, a := range values {
		newValues[i] = fmt.Sprintf("%v%%", a)
	}
	return c.Like(newValues...)
}

// IsNull operator(or fun): is null
func (c *Column[T]) IsNull(values ...any) *Predicate {
	return c.Cond(IsNullOpr, values...)
}

// IsNotNull operator(or fun): is not null
func (c *Column[T]) IsNotNull(values ...any) *Predicate {
	return c.Cond(IsNotNullOpr, values...)
}

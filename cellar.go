// Package cellar
// @author tabuyos
// @since 2023/9/11
// @description cellar
package cellar

import (
	"github.com/service-io/cellar/token"
	"strings"
)

type MRRunner func() (string, []any)
type JoinType int
type Mode int

const (
	_ JoinType = iota
	Join
	LeftJoin
	RightJoin
)

const (
	_ Mode = iota
	DftMode
	AndMode
	OrMode
)

type Selfish interface {
	Self() string
}

type Render interface {
	Render(...*strings.Builder) []any
}

type Literaler interface {
	Literal() string
}

func (j JoinType) Literal() string {
	switch j {
	case Join:
		return token.Join.Literal()
	case LeftJoin:
		return token.LeftJoin.Literal()
	case RightJoin:
		return token.RightJoin.Literal()
	default:
		return ""
	}
}

type Evaluator[T any] struct {
	defaultValues   []any
	values          []any
	sqlKey          string
	hasWhere        bool
	ref             *RefTable
	logical         *Logical
	ei              *EvalInfo[T]
	persistServices []PersistService[T]
}

func Default[T any]() *Evaluator[T] {
	return &Evaluator[T]{}
}

func WithLogicalEvaluator[T any]() *Evaluator[T] {
	return &Evaluator[T]{logical: WithLogical()}
}

func (e *Evaluator[T]) hasEnableLogical() bool {
	if e.logical == nil {
		return false
	}
	return e.logical.enable
}

func (e *Evaluator[T]) getLogicDeletedSQL() string {
	if e.hasEnableLogical() {
		if e.ref == nil {
			return e.logical.Key() + token.Equal.Pretty() + e.logical.UndeleteVal()
		}
		tables := e.ref.FlatAll()
		if len(tables) == 1 {
			return e.logical.key + token.Equal.Pretty() + e.logical.UndeleteVal()
		}
		var snips = make([]string, len(tables))
		for i, table := range tables {
			key := table.Decorate(e.logical.key)
			val := e.logical.UndeleteVal()
			snips[i] = key + token.Equal.Pretty() + val
		}
		return strings.Join(snips, token.And.Pretty())
	}
	return ""
}

func (e *Evaluator[T]) hasEvalInfo() bool {
	return e.ei != nil
}

func (e *Evaluator[T]) readEvalInfo() {
	// 遍历所有的持久化服务, 并从中获取执行信息
	for _, ps := range e.persistServices {
		if len(e.sqlKey) == 0 {
			break
		}
		lookup := ps.Lookup(e.sqlKey)
		if lookup != nil {
			e.ei = WithEvalInfo(lookup.SQL(), lookup.TotalSQL(), nil, lookup.Mappers())
		}
	}
}

func (e *Evaluator[T]) writeEvalInfo() {
	if e.ei == nil {
		return
	}
	if len(e.sqlKey) == 0 {
		return
	}
	for _, ps := range e.persistServices {
		if len(e.sqlKey) == 0 {
			break
		}
		ps.Persistence(e.sqlKey, e.ei)
	}
}

func (e *Evaluator[T]) defaultLogicalWhere(builders ...*strings.Builder) {
	if e.hasWhere {
		return
	}
	if e.hasEnableLogical() {
		deletedSQL := e.getLogicDeletedSQL()
		if len(deletedSQL) == 0 {
			return
		}
		for _, builder := range builders {
			builder.WriteString(token.Where.Pretty())
			builder.WriteString(deletedSQL)
		}
	}
}

func (e *Evaluator[T]) Cache(key string, pss ...PersistService[T]) *Evaluator[T] {
	e.sqlKey = key
	e.persistServices = pss
	e.readEvalInfo()
	return e
}

func (e *Evaluator[T]) WithValues(values ...any) {
	e.defaultValues = values
}

func (e *Evaluator[T]) WithLogicalValue(cdv ...string) *Evaluator[T] {
	if !e.hasEnableLogical() {
		// using last value
		for _, v := range cdv {
			e.logical.cdval = v
		}
	}
	return e
}

func (e *Evaluator[T]) EnableLogical() *Evaluator[T] {
	e.logical = e.logical.Enable()
	return e
}

func (e *Evaluator[T]) DisableLogical() *Evaluator[T] {
	e.logical = e.logical.Disable()
	return e
}

func (e *Evaluator[T]) EvalInfo() *EvalInfo[T] {
	if e == nil {
		return nil
	}
	return e.ei
}

func (e *Evaluator[T]) Replace(ei *EvalInfo[T]) {
	if e == nil {
		return
	}
	e.ei = ei
}

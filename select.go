// Package cellar
// @author tabuyos
// @since 2023/9/11
// @description cellar
package cellar

import (
	"github.com/service-io/cellar/token"
	"strings"
)

type SelectEvaluator[T any] struct {
	*Evaluator[T]
	pageable bool
	pageTask Task
	tasks    []Task
	hintTask Task
	mappers  []func(*T) any
}

func (e *Evaluator[T]) Select(cols ...*Column[T]) *SelectEvaluator[T] {
	se := &SelectEvaluator[T]{Evaluator: e}
	if e.hasEvalInfo() {
		return se
	}
	task := func(builders ...*strings.Builder) []any {
		execFn := func(queryBuilder *strings.Builder) {
			queryBuilder.WriteString(token.Select.Literal())
			se.hintTask.Idle(queryBuilder)
			colLit := make([]string, len(cols))
			for i, col := range cols {
				colLit[i] = col.Literal()
				se.mappers = append(se.mappers, col.Mapper())
			}
			queryBuilder.WriteString(token.Space.Literal())
			queryBuilder.WriteString(strings.Join(colLit, token.CommaSpace.Literal()))
		}
		switch len(builders) {
		case 1:
			execFn(builders[0])
		case 2:
			execFn(builders[0])
			if se.pageable {
				totalBuilder := builders[1]
				totalBuilder.WriteString(token.Select.Literal())
				totalBuilder.WriteString(token.Space.Literal())
				totalBuilder.WriteString(token.Count.Literal())
				totalBuilder.WriteString(token.LeftParentheses.Literal())
				totalBuilder.WriteString("1")
				totalBuilder.WriteString(token.RightParentheses.Literal())
			}
		default:
			panic("not found any sql builder!")
		}
		return nil
	}
	se.tasks = append(se.tasks, task)
	return se
}

func (e *SelectEvaluator[T]) Hint(ks ...token.Token) *SelectEvaluator[T] {
	if e.hasEvalInfo() {
		return e
	}
	e.hintTask = func(builders ...*strings.Builder) []any {
		snipBuilder, release := getBuilder()
		defer release()
		for _, k := range ks {
			snipBuilder.WriteString(token.Space.Literal())
			snipBuilder.WriteString(k.Literal())
		}
		for _, builder := range builders {
			builder.WriteString(snipBuilder.String())
		}
		return nil
	}
	return e
}

func (e *SelectEvaluator[T]) From(rt *RefTable) *SelectEvaluator[T] {
	if e.hasEvalInfo() {
		builder, release := getBuilder()
		defer release()
		e.values = append(e.values, rt.Render(builder)...)
		return e
	}
	e.ref = rt
	task := func(builders ...*strings.Builder) []any {
		for _, builder := range builders {
			builder.WriteString(token.Space.Literal())
			builder.WriteString(token.From.Literal())
		}
		values := rt.Render(builders...)
		e.defaultLogicalWhere(builders...)
		return values
	}
	e.tasks = append(e.tasks, task)
	return e
}

func (e *SelectEvaluator[T]) Where(pred *Predicate, values ...any) *SelectEvaluator[T] {
	if e.hasEvalInfo() {
		builder, release := getBuilder()
		if len(values) > 0 {
			e.values = append(e.values, values...)
			return e
		}
		if len(e.defaultValues) > 0 {
			return e
		}
		defer release()
		e.values = append(e.values, pred.Render(builder)...)
		return e
	}
	e.hasWhere = true
	task := func(builders ...*strings.Builder) []any {
		snipBuilder, release := getBuilder()
		defer release()
		values := pred.Render(snipBuilder)
		predSQL := snipBuilder.String()
		deletedSQL := e.getLogicDeletedSQL()
		for _, builder := range builders {
			builder.WriteString(token.Where.Pretty())
			if e.hasEnableLogical() {
				if len(deletedSQL) == 0 {
					builder.WriteString(predSQL)
					continue
				}
				if pred.Mixed() {
					builder.WriteString(token.LeftParentheses.Literal())
					builder.WriteString(predSQL)
					builder.WriteString(token.RightParentheses.Literal())
				} else {
					builder.WriteString(predSQL)
				}
				builder.WriteString(token.And.Pretty())
				builder.WriteString(deletedSQL)
			} else {
				builder.WriteString(predSQL)
			}
		}
		return values
	}
	e.tasks = append(e.tasks, task)
	return e
}

func (e *SelectEvaluator[T]) GroupBy(cols ...Selfish) *SelectEvaluator[T] {
	task := func(builders ...*strings.Builder) []any {
		snips := make([]string, len(cols))
		for i, col := range cols {
			snips[i] = col.Self()
		}
		for _, builder := range builders {
			builder.WriteString(token.GroupBy.Pretty())
			builder.WriteString(strings.Join(snips, token.CommaSpace.Literal()))
		}
		return nil
	}
	e.tasks = append(e.tasks, task)
	return e
}

func (e *SelectEvaluator[T]) Having(pred *Predicate) *SelectEvaluator[T] {
	if e.hasEvalInfo() {
		builder, release := getBuilder()
		defer release()
		e.values = append(e.values, pred.Render(builder)...)
		return e
	}
	task := func(builders ...*strings.Builder) []any {
		for _, builder := range builders {
			builder.WriteString(token.Having.Pretty())
		}
		values := pred.Render(builders...)
		return values
	}
	e.tasks = append(e.tasks, task)
	return e
}

func (e *SelectEvaluator[T]) OrderBy(orders ...*Order) *SelectEvaluator[T] {
	if e.hasEvalInfo() {
		return e
	}
	task := func(builders ...*strings.Builder) []any {
		for _, builder := range builders {
			builder.WriteString(token.OrderBy.Pretty())
			var snips = make([]string, len(orders))
			for i, order := range orders {
				snips[i] = order.Literal()
			}
			builder.WriteString(strings.Join(snips, token.CommaSpace.Literal()))
		}
		return nil
	}
	e.tasks = append(e.tasks, task)
	return e
}

func (e *SelectEvaluator[T]) Limit(limit int64) *SelectEvaluator[T] {
	if e.hasEvalInfo() {
		e.values = append(e.values, limit)
		return e
	}
	e.pageTask = func(builders ...*strings.Builder) []any {
		builder := builders[0]
		builder.WriteString(token.Limit.Pretty())
		builder.WriteString(token.Placeholder.Literal())
		return []any{limit}
	}
	return e
}

func (e *SelectEvaluator[T]) Offset(offset int64) *SelectEvaluator[T] {
	if e.hasEvalInfo() {
		e.values = append(e.values, offset)
		return e
	}
	e.pageable = true
	preTask := e.pageTask
	e.pageTask = func(builders ...*strings.Builder) []any {
		values := preTask.Idle(builders...)
		values = append(values, offset)
		builder := builders[0]
		builder.WriteString(token.Offset.Pretty())
		builder.WriteString(token.Placeholder.Literal())
		return values
	}
	return e
}

func (e *SelectEvaluator[T]) Page(limit, offset int64) *SelectEvaluator[T] {
	if e.hasEvalInfo() {
		if e.defaultValues != nil {
			e.defaultValues = append(e.defaultValues, limit, offset)
		}
		e.values = append(e.values, limit, offset)
		return e
	}
	e.pageable = true
	e.pageTask = func(builders ...*strings.Builder) []any {
		builder := builders[0]
		builder.WriteString(token.Limit.Pretty())
		builder.WriteString(token.Placeholder.Literal())
		builder.WriteString(token.Offset.Pretty())
		builder.WriteString(token.Placeholder.Literal())
		return []any{limit, offset}
	}
	return e
}

func (e *SelectEvaluator[T]) Eval() EvalInfoService[T] {
	if e.hasEvalInfo() {
		if e.defaultValues == nil {
			e.ei.values = e.values
		} else {
			e.ei.values = e.defaultValues
		}
		return e.ei
	}

	var builders []*strings.Builder
	var values []any
	execSQLBuf, execSQLRelease := getBuilder()
	defer execSQLRelease()
	totalSQLBuf, totalSQLRelease := getBuilder()
	defer totalSQLRelease()
	builders = append(builders, execSQLBuf)
	if e.pageable {
		builders = append(builders, totalSQLBuf)
	}
	for _, task := range e.tasks {
		values = append(values, task.Idle(builders...)...)
	}
	values = append(values, e.pageTask.Idle(execSQLBuf)...)

	for _, builder := range builders {
		builder.WriteString(token.Semicolon.Literal())
	}
	if e.pageable {
		e.ei = WithEvalInfo(execSQLBuf.String(), totalSQLBuf.String(), values, e.mappers)
	} else {
		e.ei = WithEvalInfo(execSQLBuf.String(), "", values, e.mappers)
	}
	e.writeEvalInfo()
	return e.ei
}

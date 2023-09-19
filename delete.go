// Package cellar
// @author tabuyos
// @since 2023/9/11
// @description cellar
package cellar

import (
	"github.com/service-io/cellar/token"
	"strings"
)

type DeleteEvaluator[T any] struct {
	*Evaluator[T]
	selfish   Selfish
	mainTask  Task
	whereTask Task
}

func (e *Evaluator[T]) Delete() *DeleteEvaluator[T] {
	return &DeleteEvaluator[T]{Evaluator: e}
}

func (e *DeleteEvaluator[T]) From(ref Selfish) *DeleteEvaluator[T] {
	if e.hasEvalInfo() {
		return e
	}
	e.selfish = ref
	if e.hasEnableLogical() {
		e.mainTask = func(builders ...*strings.Builder) []any {
			for _, builder := range builders {
				builder.WriteString(token.Update.Literal())
				builder.WriteString(token.Space.Literal())
				builder.WriteString(ref.Self())
				builder.WriteString(token.Set.Pretty())
				builder.WriteString(e.logical.key)
				builder.WriteString(token.Equal.Pretty())
				builder.WriteString(e.logical.DeletedVal())
			}
			e.defaultLogicalWhere(builders...)
			return nil
		}
	} else {
		e.PhysicalDeleted()
	}
	return e
}

func (e *DeleteEvaluator[T]) Where(pred *Predicate) *DeleteEvaluator[T] {
	if e.hasEvalInfo() {
		builder, release := getBuilder()
		defer release()
		e.values = append(e.values, pred.Render(builder)...)
		return e
	}
	e.hasWhere = true
	e.whereTask = func(builders ...*strings.Builder) []any {
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
	return e
}

func (e *DeleteEvaluator[T]) PhysicalDeleted() *DeleteEvaluator[T] {
	if e.hasEvalInfo() {
		return e
	}
	e.mainTask = func(builders ...*strings.Builder) []any {
		for _, builder := range builders {
			builder.WriteString(token.Delete.Literal())
			builder.WriteString(token.From.Pretty())
			builder.WriteString(e.selfish.Self())
		}
		e.defaultLogicalWhere(builders...)
		return nil
	}
	return e
}

func (e *DeleteEvaluator[T]) Eval() EvalInfoService[T] {
	if e.hasEvalInfo() {
		if e.defaultValues == nil {
			e.ei.values = e.values
		} else {
			e.ei.values = e.defaultValues
		}
		return e.ei
	}

	builder, release := getBuilder()
	defer release()

	e.values = append(e.values, e.mainTask.Idle(builder)...)
	e.values = append(e.values, e.whereTask.Idle(builder)...)

	builder.WriteString(token.Semicolon.Literal())

	e.ei = WithEvalInfo[T](builder.String(), "", e.values, nil)
	e.writeEvalInfo()

	return e.ei
}

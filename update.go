// Package cellar
// @author tabuyos
// @since 2023/9/11
// @description cellar
package cellar

import (
	"github.com/service-io/cellar/token"
	"strings"
)

type UpdateEvaluator[T any] struct {
	*Evaluator[T]
	mainTask  Task
	setTask   Task
	whereTask Task
	selfishs  []Selfish
}

func (e *Evaluator[T]) UpdateRef(ref Selfish, selfishs ...Selfish) *UpdateEvaluator[T] {
	ue := e.Update(ref)
	if e.hasEvalInfo() {
		return ue
	}
	ue.selfishs = selfishs
	return ue
}

func (e *UpdateEvaluator[T]) SetValues(values ...any) *UpdateEvaluator[T] {
	e.values = values
	return e
}

func (e *Evaluator[T]) Update(ref Selfish) *UpdateEvaluator[T] {
	ue := &UpdateEvaluator[T]{Evaluator: e}
	if e.hasEvalInfo() {
		return ue
	}
	ue.mainTask = func(builders ...*strings.Builder) []any {
		for _, builder := range builders {
			builder.WriteString(token.Update.Literal())
			builder.WriteString(token.Space.Literal())
			builder.WriteString(ref.Self())
		}
		return nil
	}
	ue.setTask = func(builders ...*strings.Builder) []any {
		snips := make([]string, len(ue.selfishs))
		for i, selfish := range ue.selfishs {
			snips[i] = selfish.Self() + token.Space.Join(token.EqualPlaceholder).Literal()
		}
		for _, builder := range builders {
			builder.WriteString(token.Set.Pretty())
			builder.WriteString(strings.Join(snips, token.CommaSpace.Literal()))
		}
		ue.defaultLogicalWhere(builders...)
		return nil
	}
	return ue
}

func (e *UpdateEvaluator[T]) Set(selfish Selfish, value any) *UpdateEvaluator[T] {
	e.values = append(e.values, value)
	if e.hasEvalInfo() {
		return e
	}
	e.selfishs = append(e.selfishs, selfish)
	return e
}

func (e *UpdateEvaluator[T]) Where(pred *Predicate) *UpdateEvaluator[T] {
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

func (e *UpdateEvaluator[T]) Eval() EvalInfoService[T] {
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

	e.mainTask.Idle(builder)
	e.setTask.Idle(builder)
	e.values = append(e.values, e.whereTask.Idle(builder)...)
	builder.WriteString(token.Semicolon.Literal())

	e.ei = WithEvalInfo[T](builder.String(), "", e.values, nil)
	e.writeEvalInfo()
	return e.ei
}

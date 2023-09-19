// Package cellar
// @author tabuyos
// @since 2023/9/11
// @description cellar
package cellar

import (
	"github.com/service-io/cellar/token"
	"strings"
)

type InsertEvaluator[T any] struct {
	*Evaluator[T]
	counter    int
	insertTask Task
	intoTask   Task
	valueTask  Task
}

func (e *Evaluator[T]) Insert(cols ...Selfish) *InsertEvaluator[T] {
	ie := &InsertEvaluator[T]{Evaluator: e}
	if e.hasEvalInfo() {
		return ie
	}
	ie.counter = len(cols)
	ie.insertTask = func(builders ...*strings.Builder) []any {
		snips := make([]string, len(cols))
		for i, col := range cols {
			snips[i] = col.Self()
		}
		for _, builder := range builders {
			builder.WriteString(token.LeftParentheses.Literal())
			builder.WriteString(strings.Join(snips, token.CommaSpace.Literal()))
			builder.WriteString(token.RightParentheses.Literal())
		}
		return nil
	}
	return ie
}

func (e *InsertEvaluator[T]) Into(ref Selfish) *InsertEvaluator[T] {
	if e.hasEvalInfo() {
		return e
	}
	e.intoTask = func(builders ...*strings.Builder) []any {
		for _, builder := range builders {
			builder.WriteString(token.Insert.Literal())
			builder.WriteString(token.Into.Pretty())
			builder.WriteString(ref.Self())
		}
		e.insertTask(builders...)
		return nil
	}
	return e
}

func (e *InsertEvaluator[T]) Values(values ...any) *InsertEvaluator[T] {
	e.values = append(e.values, values...)
	if e.hasEvalInfo() {
		return e
	}
	e.valueTask = func(builders ...*strings.Builder) []any {
		for _, builder := range builders {
			builder.WriteString(token.Values.Pretty())
		}
		return nil
	}
	return e
}

func (e *InsertEvaluator[T]) Value(values ...any) *InsertEvaluator[T] {
	e.values = values
	if e.hasEvalInfo() {
		return e
	}
	e.valueTask = func(builders ...*strings.Builder) []any {
		for _, builder := range builders {
			builder.WriteString(token.Value.Pretty())
		}
		return nil
	}
	return e
}

func (e *InsertEvaluator[T]) getInsertPlaceholder() string {
	var snips = make([]string, e.counter)
	for i := 0; i < e.counter; i++ {
		snips[i] = "?"
	}
	return strings.Join(snips, token.CommaSpace.Literal())
}

func (e *InsertEvaluator[T]) getTimesPlaceholder(times int, ph string) string {
	var snips = make([]string, times)
	for i := 0; i < times; i++ {
		snips[i] = token.JoinParentheses(ph)
	}
	return strings.Join(snips, ", ")
}

func (e *InsertEvaluator[T]) Eval() EvalInfoService[T] {
	if e.hasEvalInfo() {
		if e.defaultValues == nil {
			e.ei.values = e.values
		} else {
			e.ei.values = e.defaultValues
		}
		return e.ei
	}

	if e.counter == 0 {
		panic("not found any column.")
	}
	if len(e.values)%e.counter != 0 {
		panic("parameter not match.")
	}

	quotient := len(e.values) / e.counter
	placeholder := e.getInsertPlaceholder()
	timesPlaceholder := e.getTimesPlaceholder(quotient, placeholder)

	builder, release := getBuilder()
	defer release()
	e.intoTask.Idle(builder)
	e.valueTask.Idle(builder)

	builder.WriteString(timesPlaceholder)
	builder.WriteString(token.Semicolon.Literal())

	e.ei = WithEvalInfo[T](builder.String(), "", e.values, nil)
	e.writeEvalInfo()
	return e.ei
}

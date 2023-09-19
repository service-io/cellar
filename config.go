// Package cellar
// @author tabuyos
// @since 2023/9/11
// @description cellar
package cellar

type EvalService[T any] interface {
	Eval() EvalInfoService[T]
}

type EvalInfoService[T any] interface {
	EvalInfo() *EvalInfo[T]
	Replace(ei *EvalInfo[T])
}

type ConfigService[T any] interface {
	ColumnAndValue(fns ...func(*Column[T], any) bool) (selfishs []Selfish, values []any)
	Configure(func(*Evaluator[T]))
	Evaluator() *Evaluator[T]
	Asterisk(fns ...func(string) string) []*Column[T]
	PKey() *Column[T]
	LogicDelKey() *Column[T]
	Table() *RefTable
	EnableDecorate() *T
	DisableDecorate() *T
	Self() *T
}

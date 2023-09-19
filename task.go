// Package cellar
// @author tabuyos
// @since 2023/9/11
// @description cellar
package cellar

import "strings"

type Task func(...*strings.Builder) []any

func (task Task) Idle(builders ...*strings.Builder) []any {
	if task == nil {
		return nil
	}
	return task(builders...)
}

// Package cellar
// @author tabuyos
// @since 2023/9/11
// @description cellar
package cellar

import (
	"github.com/service-io/cellar/token"
)

type Order struct {
	col string
	asc bool
}

func (o *Order) Literal() string {
	if o == nil {
		return ""
	}
	if o.asc {
		return o.col
	}
	return o.col + token.Space.Join(token.Desc).Literal()
}

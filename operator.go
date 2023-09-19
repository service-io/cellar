// Package cellar
// @author tabuyos
// @since 2023/9/11
// @description cellar
package cellar

import (
	"github.com/service-io/cellar/token"
	"strings"
)

type Operator func(bool) string

func (op Operator) String() string {
	return op.Self()
}

func (op Operator) Literal() string {
	return op(true)
}

func (op Operator) Self() string {
	return op(false)
}

func AndOpr(bool) string {
	return token.And.Literal()
}

func OrOpr(bool) string {
	return token.Or.Literal()
}

// ============================== COL ==============================

func GenINOpr(count int) func(bool) string {
	if count == 1 {
		return EQOpr
	}
	var snips = make([]string, count)
	for i := 0; i < count; i++ {
		snips[i] = token.Placeholder.Literal()
	}
	return func(bool) string {
		return token.In.Join(token.LeftParentheses).JoinLit(strings.Join(snips, token.CommaSpace.Literal())).Join(token.RightParentheses).Literal()
	}
}

func EQOpr(has bool) string {
	if has {
		return "= ?"
	}
	return "="
}

func NQOpr(has bool) string {
	if has {
		return "<> ?"
	}
	return "<>"
}

func LTOpr(has bool) string {
	if has {
		return "< ?"
	}
	return "<"
}

func GTOpr(has bool) string {
	if has {
		return "> ?"
	}
	return ">"
}

func GEOpr(has bool) string {
	if has {
		return ">= ?"
	}
	return ">="
}

func LEOpr(has bool) string {
	if has {
		return "<= ?"
	}
	return "<="
}

func BetweenOpr(has bool) string {
	if has {
		var builder strings.Builder
		builder.WriteString(token.Between.Literal())
		builder.WriteString(token.SpacePlaceholder.Literal())
		builder.WriteString(token.Space.Literal())
		builder.WriteString(token.And.Literal())
		builder.WriteString(token.SpacePlaceholder.Literal())
		return builder.String()
	}
	return token.Between.Literal()
}

func NotBetweenOpr(has bool) string {
	if has {
		var builder strings.Builder
		builder.WriteString(token.Not.Literal())
		builder.WriteString(token.Space.Literal())
		builder.WriteString(token.Between.Literal())
		builder.WriteString(token.SpacePlaceholder.Literal())
		builder.WriteString(token.Space.Literal())
		builder.WriteString(token.And.Literal())
		builder.WriteString(token.SpacePlaceholder.Literal())
		return builder.String()
	}
	return token.Not.Join(token.Between).Literal()
}

func MultiOpr(bool) string {
	return "*"
}

func AddOpr(bool) string {
	return "+"
}

func MinusOpr(bool) string {
	return "-"
}

func LikeOpr(has bool) string {
	if has {
		return strings.Join([]string{token.Like.Literal(), "?"}, " ")
	}
	return token.Like.Literal()
}

func IsNullOpr(bool) string {
	return strings.Join([]string{token.Is.Literal(), token.Null.Literal()}, " ")
}

func IsNotNullOpr(bool) string {
	return strings.Join([]string{token.Is.Literal(), token.Not.Literal(), token.Null.Literal()}, " ")
}

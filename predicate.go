// Package cellar
// @author tabuyos
// @since 2023/9/11
// @description cellar
package cellar

import (
	"github.com/service-io/cellar/token"
	"strings"
)

type Predicate struct {
	mod Mode
	col Selfish
	ars []any
	lft *Predicate
	opr Operator
	rht *Predicate
}

func WithPred(self Selfish, opr Operator, ars ...any) *Predicate {
	pred := &Predicate{
		mod: DftMode,
		ars: ars,
		opr: opr,
		col: self,
	}
	return pred
}

func (p *Predicate) Mixed() bool {
	return p.mod == OrMode
}

func (p *Predicate) And(preds ...*Predicate) *Predicate {
	var fp = p
	for _, pred := range preds {
		fp = &Predicate{
			mod: AndMode,
			lft: fp,
			opr: AndOpr,
			rht: pred,
		}
	}
	return fp
}

func (p *Predicate) Or(preds ...*Predicate) *Predicate {
	var fp = p
	for _, pred := range preds {
		fp = &Predicate{
			mod: OrMode,
			lft: fp,
			opr: OrOpr,
			rht: pred,
		}
	}
	return fp
}

func (p *Predicate) Literal() (sql string, values []any) {
	builder, release := getBuilder()
	defer release()
	p.Render(builder)
	return builder.String(), p.ars
}

func (p *Predicate) equalLft(builder *strings.Builder, lftSQL, rhtSQL string) {
	if p.mod == p.rht.mod {
		builder.WriteString(lftSQL)
		builder.WriteString(token.Space.Literal())
		builder.WriteString(p.opr.Literal())
		builder.WriteString(token.Space.Literal())
		builder.WriteString(rhtSQL)
	} else {
		builder.WriteString(lftSQL)
		builder.WriteString(token.Space.Literal())
		builder.WriteString(p.opr.Literal())
		builder.WriteString(token.Space.Literal())
		builder.WriteString(token.LeftParentheses.Literal())
		builder.WriteString(rhtSQL)
		builder.WriteString(token.RightParentheses.Literal())
	}
}

func (p *Predicate) notEqualLft(builder *strings.Builder, lftSQL, rhtSQL string) {
	if p.mod == p.rht.mod {
		builder.WriteString(token.LeftParentheses.Literal())
		builder.WriteString(lftSQL)
		builder.WriteString(token.RightParentheses.Literal())
		builder.WriteString(token.Space.Literal())
		builder.WriteString(p.opr.Literal())
		builder.WriteString(token.Space.Literal())
		builder.WriteString(rhtSQL)
	} else {
		builder.WriteString(token.LeftParentheses.Literal())
		builder.WriteString(lftSQL)
		builder.WriteString(token.RightParentheses.Literal())
		builder.WriteString(token.Space.Literal())
		builder.WriteString(p.opr.Literal())
		builder.WriteString(token.Space.Literal())
		builder.WriteString(token.LeftParentheses.Literal())
		builder.WriteString(rhtSQL)
		builder.WriteString(token.RightParentheses.Literal())
	}
}

func (p *Predicate) processLftAndRht(builder *strings.Builder) {
	lftSQL, lftValues := p.lft.Literal()
	rhtSQL, rhtValues := p.rht.Literal()
	if p.lft.mod == DftMode {
		p.lft.mod = p.mod
	}
	if p.rht.mod == DftMode {
		p.rht.mod = p.mod
	}
	if p.mod == p.lft.mod {
		p.equalLft(builder, lftSQL, rhtSQL)
	} else {
		p.notEqualLft(builder, lftSQL, rhtSQL)
	}
	p.ars = append(p.ars, lftValues...)
	p.ars = append(p.ars, rhtValues...)
}

func (p *Predicate) Render(builders ...*strings.Builder) []any {
	snipBuilder, release := getBuilder()
	defer release()

	if p.lft == nil {
		snipBuilder.WriteString(p.col.Self())
		snipBuilder.WriteString(token.Space.Literal())
		snipBuilder.WriteString(p.opr.Literal())
	} else {
		p.processLftAndRht(snipBuilder)
	}

	for _, builder := range builders {
		builder.WriteString(snipBuilder.String())
	}

	return p.ars
}

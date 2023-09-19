// Package cellar
// @author tabuyos
// @since 2023/9/11
// @description cellar
package cellar

type Logical struct {
	// enable logical deleted
	enable bool
	// deleted key
	key string
	// done deleted values
	ddval string
	// undone deleted values
	udval string
	// current deleted values
	cdval string
}

func WithLogical() *Logical {
	return &Logical{
		enable: true,
		key:    "`deleted`",
		ddval:  "1",
		udval:  "0",
	}
}

func (l *Logical) CurrentVal(val string) *Logical {
	var logical = l
	if logical == nil {
		logical = WithLogical()
	}
	logical.cdval = val
	return logical
}

func (l *Logical) Key() string {
	return l.key
}

func (l *Logical) DeletedVal() string {
	if len(l.cdval) > 0 {
		return l.cdval
	}
	return l.ddval
}

func (l *Logical) UndeleteVal() string {
	if len(l.cdval) > 0 {
		return l.cdval
	}
	return l.udval
}

func (l *Logical) Enable() *Logical {
	var logical = l
	if logical == nil {
		logical = WithLogical()
	}
	logical.enable = true
	return logical
}

func (l *Logical) Disable() *Logical {
	var logical = l
	if logical == nil {
		logical = WithLogical()
	}
	logical.enable = false
	return logical
}

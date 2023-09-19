package token

import "strings"

type Token string

var lowercase bool

const (
	Space            Token = " "
	Comma            Token = ","
	CommaSpace       Token = ", "
	Semicolon        Token = ";"
	Dot              Token = "."
	Equal            Token = "="
	Placeholder      Token = "?"
	SpacePlaceholder Token = " ?"
	EqualPlaceholder Token = "= ?"
	LeftParentheses  Token = "("
	RightParentheses Token = ")"
)

const (
	As        Token = "AS"
	Join      Token = "JOIN"
	LeftJoin  Token = "LEFT JOIN"
	RightJoin Token = "RIGHT JOIN"
	Left      Token = "LEFT"
	Right     Token = "RIGHT"
	Inner     Token = "INNER"
	Outer     Token = "OUTER"
	Distinct  Token = "DISTINCT"
	From      Token = "FROM"
	Select    Token = "SELECT"
	Delete    Token = "DELETE"
	Update    Token = "UPDATE"
	Insert    Token = "INSERT"
	On        Token = "ON"
	In        Token = "IN"
	Set       Token = "SET"
	Into      Token = "INTO"
	Value     Token = "VALUE"
	Values    Token = "VALUES"
	Asterisk  Token = "ASTERISK"
	Where     Token = "WHERE"
	Between   Token = "BETWEEN"
	And       Token = "AND"
	Or        Token = "OR"
	Group     Token = "GROUP"
	GroupBy   Token = "GROUP BY"
	Having    Token = "HAVING"
	Order     Token = "ORDER"
	OrderBy   Token = "ORDER BY"
	By        Token = "BY"
	Asc       Token = "ASC"
	Desc      Token = "DESC"
	Limit     Token = "LIMIT"
	Like      Token = "Like"
	Offset    Token = "OFFSET"
	Count     Token = "COUNT"
	Is        Token = "IS"
	Not       Token = "NOT"
	Null      Token = "NULL"
)

func Of(k string) Token {
	return Token(k)
}

func LowerCase() {
	lowercase = true
}

func UpperCase() {
	lowercase = false
}

func JoinParentheses(tk string) string {
	return LeftParentheses.Literal() + tk + RightParentheses.Literal()
}

func (t Token) Pretty() string {
	return Space.Join(t).Join(Space).Literal()
}

func (t Token) String() string {
	return t.Literal()
}

func (t Token) Literal() string {
	if lowercase {
		return strings.ToLower(string(t))
	}
	return string(t)
}

func (t Token) Join(tk Token) Token {
	return t + tk
}

func (t Token) JoinLit(tk string) Token {
	return t + Token(tk)
}

func (t Token) LowerCase() Token {
	return Of(strings.ToLower(t.String()))
}

func (t Token) UpperCase() Token {
	return t
}

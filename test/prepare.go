// Package test
// @author tabuyos
// @since 2023/9/19
// @description test
package test

import (
	"github.com/service-io/cellar"
	"time"
)

type User struct {
	ID       *int64     `json:"id"`
	Name     *string    `json:"name"`
	Age      *uint8     `json:"age"`
	Birthday *time.Time `json:"birthday"`

	decorate bool
}

type Account struct {
	ID     *int64  `json:"id"`
	UserID *int64  `json:"userId"`
	RoleID *int64  `json:"roleId"`
	Name   *string `json:"name"`
	Pwd    *string `json:"pad"`
	Email  *string `json:"email"`

	decorate bool
}

type Role struct {
	ID    *int64  `json:"id"`
	Name  *string `json:"name"`
	Code  *string `json:"code"`
	Range *string `json:"range"`

	decorate bool
}

func (e *User) IDCol() *cellar.Column[User] {
	return cellar.WithColumn("`id`", func(rec *User) any {
		return &rec.ID
	}, func(key string) string {
		if !e.decorate {
			return key
		}
		return "`user`." + key
	})
}

func (e *User) NameCol() *cellar.Column[User] {
	return cellar.WithColumn("`name`", func(rec *User) any {
		return &rec.Name
	}, func(key string) string {
		if !e.decorate {
			return key
		}
		return "`user`." + key
	})
}

func (e *User) AgeCol() *cellar.Column[User] {
	return cellar.WithColumn("`age`", func(rec *User) any {
		return &rec.Age
	}, func(key string) string {
		if !e.decorate {
			return key
		}
		return "`user`." + key
	})
}

func (e *User) BirthdayCol() *cellar.Column[User] {
	return cellar.WithColumn("`birthday`", func(rec *User) any {
		return &rec.Birthday
	}, func(key string) string {
		if !e.decorate {
			return key
		}
		return "`user`." + key
	})
}

func (e *User) Table() *cellar.RefTable {
	return cellar.WithTable("`user`")
}

func (e *Account) IDCol() *cellar.Column[Account] {
	return cellar.WithColumn("`id`", func(rec *Account) any {
		return &rec.ID
	}, func(key string) string {
		if !e.decorate {
			return key
		}
		return "`account`." + key
	})
}
func (e *Account) UserIDCol() *cellar.Column[Account] {
	return cellar.WithColumn("`user_id`", func(rec *Account) any {
		return &rec.UserID
	}, func(key string) string {
		if !e.decorate {
			return key
		}
		return "`account`." + key
	})
}
func (e *Account) RoleIDCol() *cellar.Column[Account] {
	return cellar.WithColumn("`role_id`", func(rec *Account) any {
		return &rec.RoleID
	}, func(key string) string {
		if !e.decorate {
			return key
		}
		return "`account`." + key
	})
}
func (e *Account) NameCol() *cellar.Column[Account] {
	return cellar.WithColumn("`name`", func(rec *Account) any {
		return &rec.Name
	}, func(key string) string {
		if !e.decorate {
			return key
		}
		return "`account`." + key
	})
}
func (e *Account) PwdCol() *cellar.Column[Account] {
	return cellar.WithColumn("`pwd`", func(rec *Account) any {
		return &rec.Pwd
	}, func(key string) string {
		if !e.decorate {
			return key
		}
		return "`account`." + key
	})
}
func (e *Account) EmailCol() *cellar.Column[Account] {
	return cellar.WithColumn("`email`", func(rec *Account) any {
		return &rec.Email
	}, func(key string) string {
		if !e.decorate {
			return key
		}
		return "`account`." + key
	})
}

func (e *Account) Table() *cellar.RefTable {
	return cellar.WithTable("`account`")
}

func (e *Role) IDCol() *cellar.Column[Role] {
	return cellar.WithColumn("`id`", func(rec *Role) any {
		return &rec.ID
	}, func(key string) string {
		if !e.decorate {
			return key
		}
		return "`role`." + key
	})
}
func (e *Role) NameCol() *cellar.Column[Role] {
	return cellar.WithColumn("`name`", func(rec *Role) any {
		return &rec.Name
	}, func(key string) string {
		if !e.decorate {
			return key
		}
		return "`role`." + key
	})
}
func (e *Role) CodeCol() *cellar.Column[Role] {
	return cellar.WithColumn("`code`", func(rec *Role) any {
		return &rec.Code
	}, func(key string) string {
		if !e.decorate {
			return key
		}
		return "`role`." + key
	})
}
func (e *Role) RangeCol() *cellar.Column[Role] {
	return cellar.WithColumn("`range`", func(rec *Role) any {
		return &rec.Range
	}, func(key string) string {
		if !e.decorate {
			return key
		}
		return "`role`." + key
	})
}

func (e *Role) Table() *cellar.RefTable {
	return cellar.WithTable("`role`")
}

func (e *User) EnableDecorate() *User {
	e.decorate = true
	return e
}
func (e *User) DisableDecorate() *User {
	e.decorate = true
	return e
}

func (e *Account) EnableDecorate() *Account {
	e.decorate = true
	return e
}
func (e *Account) DisableDecorate() *Account {
	e.decorate = true
	return e
}

func (e *Role) EnableDecorate() *Role {
	e.decorate = true
	return e
}
func (e *Role) DisableDecorate() *Role {
	e.decorate = true
	return e
}

func NewUser() *User {
	return &User{}
}

func NewAccount() *Account {
	return &Account{}
}

func NewRole() *Role {
	return &Role{}
}

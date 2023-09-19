// Package test
// @author tabuyos
// @since 2023/9/11
// @description cellar test
package test

import (
	"github.com/service-io/cellar"
	"github.com/service-io/cellar/token"
	"testing"
)

func TestSelectEvaluator_Eval(t *testing.T) {
	token.LowerCase()
	user := NewUser().EnableDecorate()
	account := NewAccount().EnableDecorate()
	role := NewRole().EnableDecorate()

	t.Run("logical", func(t *testing.T) {
		eval := cellar.WithLogicalEvaluator[User]()
		predicate := user.NameCol().EQ("tabuyos").And(account.NameCol().EQ("tabuyos"), account.EmailCol().EQ("example@example.com")).Or(role.NameCol().EQ("test"), role.RangeCol().EQ("read"))
		accountRefTable := account.Table().LeftJoin().OnEQ(user.IDCol(), account.UserIDCol())
		roleRefTable := role.Table().LeftJoin().OnEQ(role.IDCol(), account.RoleIDCol())
		table := user.Table().Ref(accountRefTable, roleRefTable)
		eval.Select(user.IDCol(), user.NameCol(), user.AgeCol(), user.BirthdayCol()).Hint(token.Distinct).From(table).Where(predicate).GroupBy(user.NameCol()).Having(user.NameCol().EQ("tabuyos")).OrderBy(user.BirthdayCol().Asc(), user.AgeCol().Desc()).Limit(20).Offset(0).Eval()
		info := eval.EvalInfo()
		if len(info.Values()) != 8 {
			t.Error("values not matched...")
		}
		if info.SQL() != "select distinct `user`.`id`, `user`.`name`, `user`.`age`, `user`.`birthday` from `user` left join `account` on `user`.`id` = `account`.`user_id` left join `role` on `role`.`id` = `account`.`role_id` where ((`user`.`name` = ? and `account`.`name` = ? and `account`.`email` = ?) or `role`.`name` = ? or `role`.`range` = ?) and `user`.`deleted` = 0 and `account`.`deleted` = 0 and `role`.`deleted` = 0 group by `user`.`name` having `user`.`name` = ? order by `user`.`birthday`, `user`.`age` desc limit ? offset ?;" {
			t.Error("sql not matched...")
		}
	})

	t.Run("non-logical", func(t *testing.T) {
		eval := cellar.Default[User]()
		predicate := user.NameCol().EQ("tabuyos").And(account.NameCol().EQ("tabuyos"), account.EmailCol().EQ("example@example.com")).Or(role.NameCol().EQ("test"), role.RangeCol().EQ("read"))
		accountRefTable := account.Table().LeftJoin().OnEQ(user.IDCol(), account.UserIDCol())
		roleRefTable := role.Table().LeftJoin().OnEQ(role.IDCol(), account.RoleIDCol())
		table := user.Table().Ref(accountRefTable, roleRefTable)
		eval.Select(user.IDCol(), user.NameCol(), user.AgeCol(), user.BirthdayCol()).Hint(token.Distinct).From(table).Where(predicate).GroupBy(user.NameCol()).Having(user.NameCol().EQ("tabuyos")).OrderBy(user.BirthdayCol().Asc(), user.AgeCol().Desc()).Limit(20).Offset(0).Eval()
		info := eval.EvalInfo()
		if len(info.Values()) != 8 {
			t.Error("values not matched...")
		}
		if info.SQL() != "select distinct `user`.`id`, `user`.`name`, `user`.`age`, `user`.`birthday` from `user` left join `account` on `user`.`id` = `account`.`user_id` left join `role` on `role`.`id` = `account`.`role_id` where (`user`.`name` = ? and `account`.`name` = ? and `account`.`email` = ?) or `role`.`name` = ? or `role`.`range` = ? group by `user`.`name` having `user`.`name` = ? order by `user`.`birthday`, `user`.`age` desc limit ? offset ?;" {
			t.Error("sql not matched...")
		}
	})

}

func BenchmarkSelectEvaluator_Eval(b *testing.B) {
	token.LowerCase()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := NewUser().EnableDecorate()
		account := NewAccount().EnableDecorate()
		role := NewRole().EnableDecorate()

		eval := cellar.WithLogicalEvaluator[User]()
		predicate := user.NameCol().EQ("tabuyos").And(account.NameCol().EQ("tabuyos"), account.EmailCol().EQ("example@example.com")).Or(role.NameCol().EQ("test"), role.RangeCol().EQ("read"))
		accountRefTable := account.Table().LeftJoin().OnEQ(user.IDCol(), account.UserIDCol())
		roleRefTable := role.Table().LeftJoin().OnEQ(role.IDCol(), account.RoleIDCol())
		table := user.Table().Ref(accountRefTable, roleRefTable)
		eval.Select(user.IDCol(), user.NameCol(), user.AgeCol(), user.BirthdayCol()).Hint(token.Distinct).From(table).Where(predicate).GroupBy(user.NameCol()).Having(user.NameCol().EQ("tabuyos")).OrderBy(user.BirthdayCol().Asc(), user.AgeCol().Desc()).Limit(20).Offset(0).Eval()
	}
}

func BenchmarkParallelSelectEvaluator_Eval(b *testing.B) {
	token.LowerCase()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			user := NewUser().EnableDecorate()
			account := NewAccount().EnableDecorate()
			role := NewRole().EnableDecorate()

			eval := cellar.WithLogicalEvaluator[User]()
			predicate := user.NameCol().EQ("tabuyos").And(account.NameCol().EQ("tabuyos"), account.EmailCol().EQ("example@example.com")).Or(role.NameCol().EQ("test"), role.RangeCol().EQ("read"))
			accountRefTable := account.Table().LeftJoin().OnEQ(user.IDCol(), account.UserIDCol())
			roleRefTable := role.Table().LeftJoin().OnEQ(role.IDCol(), account.RoleIDCol())
			table := user.Table().Ref(accountRefTable, roleRefTable)
			eval.WithValues(1, 2, 3)
			eval.Select(user.IDCol(), user.NameCol(), user.AgeCol(), user.BirthdayCol()).Hint(token.Distinct).From(table).Where(predicate).GroupBy(user.NameCol()).Having(user.NameCol().EQ("tabuyos")).OrderBy(user.BirthdayCol().Asc(), user.AgeCol().Desc()).Limit(20).Offset(0).Eval()
		}
	})
}

func BenchmarkParallelPersistSelectEvaluator_Eval(b *testing.B) {
	token.LowerCase()
	sqlMap := map[string]*cellar.EvalInfo[User]{}
	memoryPersist := cellar.NewMemoryPersist[User](&lks{m: sqlMap})

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			user := NewUser().EnableDecorate()
			account := NewAccount().EnableDecorate()
			role := NewRole().EnableDecorate()

			eval := cellar.WithLogicalEvaluator[User]()
			predicate := user.NameCol().EQ("tabuyos").And(account.NameCol().EQ("tabuyos"), account.EmailCol().EQ("example@example.com")).Or(role.NameCol().EQ("test"), role.RangeCol().EQ("read"))
			accountRefTable := account.Table().LeftJoin().OnEQ(user.IDCol(), account.UserIDCol())
			roleRefTable := role.Table().LeftJoin().OnEQ(role.IDCol(), account.RoleIDCol())
			table := user.Table().Ref(accountRefTable, roleRefTable)
			eval.WithValues(1, 2, 3)
			eval.Cache("select", memoryPersist).Select(user.IDCol(), user.NameCol(), user.AgeCol(), user.BirthdayCol()).Hint(token.Distinct).From(table).Where(predicate).GroupBy(user.NameCol()).Having(user.NameCol().EQ("tabuyos")).OrderBy(user.BirthdayCol().Asc(), user.AgeCol().Desc()).Limit(20).Offset(0).Eval()
		}
	})
}

type lks struct {
	m map[string]*cellar.EvalInfo[User]
}

func (l *lks) Put(key string, info *cellar.EvalInfo[User]) {
	l.m[key] = info
}

func (l *lks) Get(key string) *cellar.EvalInfo[User] {
	return l.m[key]
}

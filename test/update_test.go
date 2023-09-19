// Package test
// @author tabuyos
// @since 2023/9/11
// @description cellar test
package test

import (
	"fmt"
	"github.com/service-io/cellar"
	"github.com/service-io/cellar/token"
	"testing"
)

func TestUpdateEvaluator_Eval(t *testing.T) {
	token.LowerCase()
	user := NewUser()

	t.Run("logical", func(t *testing.T) {
		eval := cellar.WithLogicalEvaluator[User]()
		eval.Update(user.Table()).Set(user.NameCol(), 321).Set(user.AgeCol(), 26).Set(user.BirthdayCol(), 1).Where(user.IDCol().EQ("111000")).Eval()
		info := eval.EvalInfo()
		fmt.Println(info.SQL())
		if len(info.Values()) != 4 {
			t.Error("values not matched...")
		}
		if info.SQL() != "update `user` set `name` = ?, `age` = ?, `birthday` = ? where `id` = ? and `deleted` = 0;" {
			t.Error("sql not matched...")
		}
	})

	t.Run("non-logical", func(t *testing.T) {
		eval := cellar.Default[User]()
		eval.Update(user.Table()).Set(user.NameCol(), 321).Set(user.AgeCol(), 26).Set(user.BirthdayCol(), 1).Where(user.IDCol().EQ("111000")).Eval()
		info := eval.EvalInfo()
		if len(info.Values()) != 4 {
			t.Error("values not matched...")
		}
		if info.SQL() != "update `user` set `name` = ?, `age` = ?, `birthday` = ? where `id` = ?;" {
			t.Error("sql not matched...")
		}
	})

}

func BenchmarkUpdateEvaluator_Eval(b *testing.B) {
	token.LowerCase()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := NewUser()

		eval := cellar.WithLogicalEvaluator[User]()
		eval.Update(user.Table()).Set(user.NameCol(), 321).Set(user.AgeCol(), 26).Set(user.BirthdayCol(), 1).Where(user.IDCol().EQ("111000")).Eval()
	}
}

func BenchmarkParallelUpdateEvaluator_Eval(b *testing.B) {
	token.LowerCase()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			user := NewUser()

			eval := cellar.WithLogicalEvaluator[User]()
			eval.Update(user.Table()).Set(user.NameCol(), 321).Set(user.AgeCol(), 26).Set(user.BirthdayCol(), 1).Where(user.IDCol().EQ("111000")).Eval()
		}
	})
}

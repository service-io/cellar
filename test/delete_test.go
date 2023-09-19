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

func TestDeleteEvaluator_Eval(t *testing.T) {
	token.LowerCase()
	user := NewUser()

	t.Run("logical", func(t *testing.T) {
		eval := cellar.WithLogicalEvaluator[User]()
		eval.Delete().From(user.Table()).Where(user.NameCol().EQ(12).Or(user.AgeCol().EQ(24).And(user.IDCol().EQ("111000")))).Eval()
		info := eval.EvalInfo()
		if len(info.Values()) != 3 {
			t.Error("values not matched...")
		}
		if info.SQL() != "update `user` set `deleted` = 1 where (`name` = ? or (`age` = ? and `id` = ?)) and `deleted` = 0;" {
			t.Error("sql not matched...")
		}
	})

	t.Run("non-logical", func(t *testing.T) {
		eval := cellar.Default[User]()
		eval.Delete().From(user.Table()).Where(user.NameCol().EQ(12).Or(user.AgeCol().EQ(24).And(user.IDCol().EQ("111000")))).Eval()
		info := eval.EvalInfo()
		if len(info.Values()) != 3 {
			t.Error("values not matched...")
		}
		if info.SQL() != "delete from `user` where `name` = ? or (`age` = ? and `id` = ?);" {
			t.Error("sql not matched...")
		}
	})

}

func BenchmarkDeleteEvaluator_Eval(b *testing.B) {
	token.LowerCase()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := NewUser()

		eval := cellar.WithLogicalEvaluator[User]()
		eval.Delete().From(user.Table()).Where(user.NameCol().EQ(12).Or(user.AgeCol().EQ(24).And(user.IDCol().EQ("111000")))).Eval()
	}
}

func BenchmarkParallelDeleteEvaluator_Eval(b *testing.B) {
	token.LowerCase()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			user := NewUser()

			eval := cellar.WithLogicalEvaluator[User]()
			eval.Delete().From(user.Table()).Where(user.NameCol().EQ(12).Or(user.AgeCol().EQ(24).And(user.IDCol().EQ("111000")))).Eval()
		}
	})
}

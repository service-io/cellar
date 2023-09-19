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

func TestInsertEvaluator_Eval(t *testing.T) {
	token.LowerCase()
	user := NewUser()

	eval := cellar.WithLogicalEvaluator[User]()
	eval.Insert(user.NameCol(), user.AgeCol(), user.BirthdayCol()).Into(user.Table()).Values("1", "2", "3", "1", "2", "3", "1", "2", "3").Eval()
	info := eval.EvalInfo()

	if len(info.Values()) != 9 {
		t.Error("values not matched...")
	}
	if info.SQL() != "insert into `user`(`name`, `age`, `birthday`) values (?, ?, ?), (?, ?, ?), (?, ?, ?);" {
		t.Error("sql not matched...")
	}
}

func BenchmarkInsertEvaluator_Eval(b *testing.B) {
	token.LowerCase()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := NewUser()

		eval := cellar.WithLogicalEvaluator[User]()
		eval.Insert(user.NameCol(), user.AgeCol(), user.BirthdayCol()).Into(user.Table()).Values("1", "2", "3", "1", "2", "3", "1", "2", "3").Eval()
	}
}

func BenchmarkParallelInsertEvaluator_Eval(b *testing.B) {
	token.LowerCase()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			user := NewUser()

			eval := cellar.WithLogicalEvaluator[User]()
			eval.Insert(user.NameCol(), user.AgeCol(), user.BirthdayCol()).Into(user.Table()).Values("1", "2", "3", "1", "2", "3", "1", "2", "3").Eval()
		}
	})
}

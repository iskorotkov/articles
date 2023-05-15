package fizzbuzz_test

import (
	"context"
	"testing"

	"articles/src/fizzbuzz"
)

func BenchmarkValue(b *testing.B) {
	var (
		req  = &fizzbuzz.ControllerReq{From: 1, To: 100}
		resp *fizzbuzz.ControllerResp
		err  error
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err = fizzbuzz.ValueController(context.TODO(), req)
	}

	_, _ = resp, err
}

func BenchmarkPointer(b *testing.B) {
	var (
		req  = &fizzbuzz.ControllerReq{From: 1, To: 100}
		resp *fizzbuzz.ControllerResp
		err  error
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err = fizzbuzz.PtrController(context.TODO(), req)
	}

	_, _ = resp, err
}

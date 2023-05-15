package pointers_test

import (
	"testing"

	"articles/src/pointers"
)

func BenchmarkMapValues8KB(b *testing.B) {
	var (
		user = createUser8KB()
		res  pointers.User8KB
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dto := pointers.UserToDTO8KB(user)
		res = pointers.DTOToUser8KB(dto)
	}

	_ = res
}

func BenchmarkMapPointers8KB(b *testing.B) {
	var (
		user = createUser8KB()
		res  *pointers.User8KB
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dto := pointers.UserPtrToDTO8KB(&user)
		res = pointers.DTOPtrToUser8KB(dto)
	}

	_ = res
}

func createUser8KB() pointers.User8KB {
	return pointers.User8KB{
		Data: [8192]byte{},
	}
}

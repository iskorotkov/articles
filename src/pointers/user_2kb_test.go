package pointers_test

import (
	"testing"

	"articles/src/pointers"
)

func BenchmarkMapValues2KB(b *testing.B) {
	var (
		user = createUser2KB()
		res  pointers.User2KB
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dto := pointers.UserToDTO2KB(user)
		res = pointers.DTOToUser2KB(dto)
	}

	_ = res
}

func BenchmarkMapPointers2KB(b *testing.B) {
	var (
		user = createUser2KB()
		res  *pointers.User2KB
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dto := pointers.UserPtrToDTO2KB(&user)
		res = pointers.DTOPtrToUser2KB(dto)
	}

	_ = res
}

func createUser2KB() pointers.User2KB {
	return pointers.User2KB{
		Data: [2048]byte{},
	}
}

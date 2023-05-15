package pointers_test

import (
	"testing"

	"articles/src/pointers"
)

func BenchmarkMapValues1MB(b *testing.B) {
	var (
		user = createUser1MB()
		res  pointers.User1MB
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dto := pointers.UserToDTO1MB(user)
		res = pointers.DTOToUser1MB(dto)
	}

	_ = res
}

func BenchmarkMapPointers1MB(b *testing.B) {
	var (
		user = createUser1MB()
		res  *pointers.User1MB
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dto := pointers.UserPtrToDTO1MB(&user)
		res = pointers.DTOPtrToUser1MB(dto)
	}

	_ = res
}

func createUser1MB() pointers.User1MB {
	return pointers.User1MB{
		Data: [1024*1024]byte{},
	}
}

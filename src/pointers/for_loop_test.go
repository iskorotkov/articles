package pointers_test

import (
	"testing"

	"articles/src/pointers"
)

func BenchmarkMapSliceOfValues1MB(b *testing.B) {
	var (
		users = createSliceOfUser1MB()
		res   []pointers.UserDTO1MB
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res = pointers.UsersToDTOs(users)
	}

	_ = res
}

func BenchmarkMapSliceOfPointers1MB(b *testing.B) {
	var (
		users = createSliceOfPointersToUser1MB()
		res   []*pointers.UserDTO1MB
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res = pointers.UsersPtrsToDTOs(users)
	}

	_ = res
}

func createSliceOfUser1MB() []pointers.User1MB {
	const size = 10

	res := make([]pointers.User1MB, 0, size)
	for i := 0; i < size; i++ {
		res = append(res, createUser1MB())
	}

	return res
}

func createSliceOfPointersToUser1MB() []*pointers.User1MB {
	const size = 10

	res := make([]*pointers.User1MB, 0, size)
	for i := 0; i < size; i++ {
		user := createUser1MB()
		res = append(res, &user)
	}

	return res
}

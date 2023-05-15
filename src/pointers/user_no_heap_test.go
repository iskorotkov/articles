package pointers_test

import (
	"math/big"
	"testing"
	"time"

	"articles/src/pointers"
)

func BenchmarkMapValuesNoHeap(b *testing.B) {
	var (
		user = createUserNoHeap()
		res  pointers.UserNoHeap
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dto := pointers.UserToDTONoHeap(user)
		res = pointers.DTOToUserNoHeap(dto)
	}

	_ = res
}

func BenchmarkMapPointersNoHeap(b *testing.B) {
	var (
		user = createUserNoHeap()
		res  *pointers.UserNoHeap
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dto := pointers.UserPtrToDTONoHeap(&user)
		res = pointers.DTOPtrToUserNoHeap(dto)
	}

	_ = res
}

func createUserNoHeap() pointers.UserNoHeap {
	return pointers.UserNoHeap{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Now(),

		FirstName:   "John",
		SecondName:  "Doe",
		Patronymic:  "Smith",
		Birthday:    time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		Nationality: "Russian",
		UserType:    1,

		Balance:     big.NewRat(1000, 1),
		BonusPoints: big.NewRat(100, 1),
	}
}

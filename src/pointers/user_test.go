package pointers_test

import (
	"math/big"
	"testing"
	"time"

	"articles/src/pointers"
)

func BenchmarkMapValues(b *testing.B) {
	var (
		user = createUser()
		res  pointers.User
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dto := pointers.UserToDTO(user)
		res = pointers.DTOToUser(dto)
	}

	_ = res
}

func BenchmarkMapPointers(b *testing.B) {
	var (
		user = createUser()
		res  *pointers.User
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dto := pointers.UserPtrToDTO(&user)
		res = pointers.DTOPtrToUser(dto)
	}

	_ = res
}

func createUser() pointers.User {
	return pointers.User{
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

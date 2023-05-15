package pointers

import (
	"math/big"
	"time"
)

type UserNoHeap struct {
	ID                              int64
	CreatedAt, UpdatedAt, DeletedAt time.Time

	FirstName, SecondName, Patronymic string
	Birthday                          time.Time
	Nationality                       string
	UserType                          int

	Balance     *big.Rat
	BonusPoints *big.Rat
}

type UserDTONoHeap struct {
	ID                              int64
	CreatedAt, UpdatedAt, DeletedAt time.Time

	FirstName, SecondName, Patronymic string
	Birthday                          time.Time
	Nationality                       string
	UserType                          int

	Balance     *big.Rat
	BonusPoints *big.Rat
}

func UserToDTONoHeap(u UserNoHeap) UserDTONoHeap {
	return UserDTONoHeap{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,

		FirstName:   u.FirstName,
		SecondName:  u.SecondName,
		Patronymic:  u.Patronymic,
		Birthday:    u.Birthday,
		Nationality: u.Nationality,
		UserType:    u.UserType,

		Balance:     u.Balance,
		BonusPoints: u.BonusPoints,
	}
}

func DTOToUserNoHeap(d UserDTONoHeap) UserNoHeap {
	return UserNoHeap{
		ID:        d.ID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		DeletedAt: d.DeletedAt,

		FirstName:   d.FirstName,
		SecondName:  d.SecondName,
		Patronymic:  d.Patronymic,
		Birthday:    d.Birthday,
		Nationality: d.Nationality,
		UserType:    d.UserType,

		Balance:     d.Balance,
		BonusPoints: d.BonusPoints,
	}
}

func UserPtrToDTONoHeap(u *UserNoHeap) *UserDTONoHeap {
	return &UserDTONoHeap{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,

		FirstName:   u.FirstName,
		SecondName:  u.SecondName,
		Patronymic:  u.Patronymic,
		Birthday:    u.Birthday,
		Nationality: u.Nationality,
		UserType:    u.UserType,

		Balance:     u.Balance,
		BonusPoints: u.BonusPoints,
	}
}

func DTOPtrToUserNoHeap(d *UserDTONoHeap) *UserNoHeap {
	return &UserNoHeap{
		ID:        d.ID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		DeletedAt: d.DeletedAt,

		FirstName:   d.FirstName,
		SecondName:  d.SecondName,
		Patronymic:  d.Patronymic,
		Birthday:    d.Birthday,
		Nationality: d.Nationality,
		UserType:    d.UserType,

		Balance:     d.Balance,
		BonusPoints: d.BonusPoints,
	}
}

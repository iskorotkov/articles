package pointers

import (
	"math/big"
	"time"
)

type User struct {
	ID                              int64
	CreatedAt, UpdatedAt, DeletedAt time.Time

	FirstName, SecondName, Patronymic string
	Birthday                          time.Time
	Nationality                       string
	UserType                          int

	Balance     *big.Rat
	BonusPoints *big.Rat
}

type UserDTO struct {
	ID                              int64
	CreatedAt, UpdatedAt, DeletedAt time.Time

	FirstName, SecondName, Patronymic string
	Birthday                          time.Time
	Nationality                       string
	UserType                          int

	Balance     *big.Rat
	BonusPoints *big.Rat

	FullName               string
	BalanceWithBonusPoints *big.Rat
}

func UserToDTO(u User) UserDTO {
	return UserDTO{
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

		FullName:               u.FirstName + " " + u.SecondName + " " + u.Patronymic,
		BalanceWithBonusPoints: new(big.Rat).Add(u.Balance, u.BonusPoints),
	}
}

func DTOToUser(d UserDTO) User {
	return User{
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

func UserPtrToDTO(u *User) *UserDTO {
	return &UserDTO{
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

		FullName:               u.FirstName + " " + u.SecondName + " " + u.Patronymic,
		BalanceWithBonusPoints: new(big.Rat).Add(u.Balance, u.BonusPoints),
	}
}

func DTOPtrToUser(d *UserDTO) *User {
	return &User{
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

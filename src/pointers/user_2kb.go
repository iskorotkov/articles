package pointers

type User2KB struct {
	Data [2048]byte
}

type UserDTO2KB struct {
	Data [2048]byte
}

func UserToDTO2KB(u User2KB) UserDTO2KB {
	return UserDTO2KB{
		Data: u.Data,
	}
}

func DTOToUser2KB(d UserDTO2KB) User2KB {
	return User2KB{
		Data: d.Data,
	}
}

func UserPtrToDTO2KB(u *User2KB) *UserDTO2KB {
	return &UserDTO2KB{
		Data: u.Data,
	}
}

func DTOPtrToUser2KB(d *UserDTO2KB) *User2KB {
	return &User2KB{
		Data: d.Data,
	}
}

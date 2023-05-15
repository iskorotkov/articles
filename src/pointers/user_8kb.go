package pointers

type User8KB struct {
	Data [8192]byte
}

type UserDTO8KB struct {
	Data [8192]byte
}

func UserToDTO8KB(u User8KB) UserDTO8KB {
	return UserDTO8KB{
		Data: u.Data,
	}
}

func DTOToUser8KB(d UserDTO8KB) User8KB {
	return User8KB{
		Data: d.Data,
	}
}

func UserPtrToDTO8KB(u *User8KB) *UserDTO8KB {
	return &UserDTO8KB{
		Data: u.Data,
	}
}

func DTOPtrToUser8KB(d *UserDTO8KB) *User8KB {
	return &User8KB{
		Data: d.Data,
	}
}

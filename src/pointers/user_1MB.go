package pointers

type User1MB struct {
	Data [1024 * 1024]byte
}

type UserDTO1MB struct {
	Data [1024 * 1024]byte
}

func UserToDTO1MB(u User1MB) UserDTO1MB {
	return UserDTO1MB{
		Data: u.Data,
	}
}

func DTOToUser1MB(d UserDTO1MB) User1MB {
	return User1MB{
		Data: d.Data,
	}
}

func UserPtrToDTO1MB(u *User1MB) *UserDTO1MB {
	return &UserDTO1MB{
		Data: u.Data,
	}
}

func DTOPtrToUser1MB(d *UserDTO1MB) *User1MB {
	return &User1MB{
		Data: d.Data,
	}
}

package pointers

func UsersToDTOs(users []User1MB) []UserDTO1MB {
	res := make([]UserDTO1MB, 0, len(users))
	for _, u := range users {
		res = append(res, UserToDTO1MB(u))
	}

	return res
}

func UsersPtrsToDTOs(users []*User1MB) []*UserDTO1MB {
	res := make([]*UserDTO1MB, 0, len(users))
	for _, u := range users {
		res = append(res, UserPtrToDTO1MB(u))
	}

	return res
}

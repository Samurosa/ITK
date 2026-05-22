package storage

import (
	pb "ITK_Code/m/v2/proto"
)

type UserData struct {
	Data map[string]User
}

type User struct {
	ID              string
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            pb.Role
}

func NewUserStorage() *UserData {
	return &UserData{Data: make(map[string]User)}
}

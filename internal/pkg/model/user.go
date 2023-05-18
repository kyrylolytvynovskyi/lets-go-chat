package model

import (
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID
	UserName string
	Password string
}

type CreateUserRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	Id       string `json:"id"`
	UserName string `json:"userName"`
}

func (r *CreateUserRequest) Validate() error {
	lenUserName := len(r.UserName)
	if lenUserName < 4 {
		return fmt.Errorf("userName minLength is 4, actual length is %v", lenUserName)
	}

	lenPassword := len(r.Password)
	if lenPassword < 8 {
		return fmt.Errorf("password minLength is 8, actual length is %v", lenPassword)
	}

	return nil
}

package service

import (
	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/model"
)

type User interface {
	CreateUser(req model.CreateUserRequest) (model.CreateUserResponse, error)
	LoginUser(req model.LoginUserRequest) (model.LoginUserResponse, error)

	//creational patterns: prototype
	Clone() User
}

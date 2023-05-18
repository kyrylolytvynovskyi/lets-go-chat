package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/model"
)

type User struct {
	users map[string]model.User
}

func NewUser() *User {
	return &User{users: map[string]model.User{}}
}

func (srv *User) CreateUser(req model.CreateUserRequest) (model.User, error) {
	_, exist := srv.users[req.UserName]
	if exist {
		return model.User{}, fmt.Errorf("User %s already exist", req.UserName)
	}

	newUser := model.User{Id: uuid.New(), UserName: req.UserName, Password: req.Password}
	srv.users[req.UserName] = newUser

	return newUser, nil
}

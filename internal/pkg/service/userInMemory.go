package service

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/model"
)

type UserInMemory struct {
	users map[string]model.User
}

func NewUserInMemory() User {
	return &UserInMemory{users: map[string]model.User{}}
}

func (srv *UserInMemory) CreateUser(req model.CreateUserRequest) (model.CreateUserResponse, error) {
	_, exist := srv.users[req.UserName]
	if exist {
		return model.CreateUserResponse{}, fmt.Errorf("User %s already exist", req.UserName)
	}

	newUser := model.User{Id: uuid.New(), UserName: req.UserName, Password: req.Password}
	srv.users[req.UserName] = newUser

	return model.CreateUserResponse{Id: newUser.Id.String(), UserName: newUser.UserName}, nil
}

func (srv *UserInMemory) LoginUser(req model.LoginUserRequest) (model.LoginUserResponse, error) {
	user, exist := srv.users[req.UserName]

	if !exist {
		return model.LoginUserResponse{}, fmt.Errorf("login %s invalid", req.UserName)
	}

	if user.Password != req.Password {
		return model.LoginUserResponse{}, fmt.Errorf("invalid password")
	}

	url := "ws://fancy-chat.io/ws&token=one-time-token"
	return model.LoginUserResponse{Url: url}, nil
}

func (srv *UserInMemory) Clone() User {
	users := make(map[string]model.User, len(srv.users))
	for k, v := range srv.users {
		users[k] = v
	}

	userInMemory := &UserInMemory{users}
	return userInMemory
}

// behavioral patterns: iterator
func (srv *UserInMemory) GetIterator() UserIterator {

	keys := make([]string, 0, len(srv.users))
	for k := range srv.users {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return UserIterator{userInMemory: srv, keys: keys}
}

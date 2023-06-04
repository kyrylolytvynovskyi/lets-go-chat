package service

import (
	"testing"

	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/model"
)

func TestUserInMemory_Clone(t *testing.T) {
	userInMemory := UserInMemory{users: map[string]model.User{}}
	userInMemory.CreateUser(model.CreateUserRequest{UserName: "username", Password: "pwd"})

	cloned := userInMemory.Clone()

	_, err := cloned.LoginUser(model.LoginUserRequest{UserName: "username", Password: "pwd"})
	if err != nil {
		t.Fatal("loginUser unexpectedly failed")
	}

	_, err = cloned.LoginUser(model.LoginUserRequest{UserName: "username", Password: "pwd111"})
	if err == nil {
		t.Fatal("loginUser with wrong password should not succeed")
	}
}

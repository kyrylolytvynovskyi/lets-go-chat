package service

import (
	"testing"

	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/model"
	"github.com/stretchr/testify/assert"
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

func TestUserInMemory_CreateUser(t *testing.T) {

	srv := NewUserInMemory()

	newUser := model.CreateUserRequest{UserName: "user", Password: "password"}

	//create new user
	resp, err := srv.CreateUser(newUser)
	assert.NoError(t, err)
	assert.Equal(t, "user", resp.UserName)
	assert.NotEqual(t, "", resp.Id)

	//create existing user
	resp, err = srv.CreateUser(newUser)
	assert.Error(t, err)
	assert.Equal(t, "", resp.UserName)
	assert.Equal(t, "", resp.Id)
}

func TestUserInMemory_LoginUser(t *testing.T) {

	srv := NewUserInMemory()
	newUser := model.CreateUserRequest{UserName: "user", Password: "password"}

	//create new user
	srv.CreateUser(newUser)

	//invalid login
	loginUser := model.LoginUserRequest{UserName: "user1", Password: "password"}
	_, err := srv.LoginUser(loginUser)
	assert.ErrorContains(t, err, "invalid", "Invalid login")

	//invalid password
	loginUser = model.LoginUserRequest{UserName: "user", Password: "password1"}
	_, err = srv.LoginUser(loginUser)
	assert.ErrorContains(t, err, "invalid", "Invalid password")

	//valid login
	loginUser = model.LoginUserRequest{UserName: "user", Password: "password"}
	_, err = srv.LoginUser(loginUser)
	assert.NoError(t, err)

}

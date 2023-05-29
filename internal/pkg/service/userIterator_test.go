package service

import (
	"testing"

	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/model"
)

func TestUserIterator_MoveNext_GetValue(t *testing.T) {

	userInMemory := UserInMemory{users: map[string]model.User{}}
	userInMemory.CreateUser(model.CreateUserRequest{UserName: "username1", Password: "pwd"})
	userInMemory.CreateUser(model.CreateUserRequest{UserName: "username2", Password: "pwd"})

	it := userInMemory.GetIterator()
	user, err := it.GetValue()
	if err != nil {
		t.Fatalf("user not returned")
	}

	if user.UserName != "username1" {
		t.Fatalf("wrong user received, expected username1, got %s", user.UserName)
	}

	it.MoveNext()
	user, err = it.GetValue()
	if err != nil {
		t.Fatalf("user not returned")
	}

	if user.UserName != "username2" {
		t.Fatalf("wrong user received, expected username2, got %s", user.UserName)
	}

	it.MoveNext()
	user, err = it.GetValue()
	if err == nil {
		t.Fatalf("end of iterator expected but got an object")
	}
}

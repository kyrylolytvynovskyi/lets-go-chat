package service

import (
	"fmt"

	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/model"
)

// behavioral pattern: iterator
type UserIterator struct {
	userInMemory *UserInMemory
	keys         []string
	pos          int
}

func (it *UserIterator) MoveNext() {
	it.pos++
}

func (it *UserIterator) GetValue() (model.User, error) {

	if it.pos >= len(it.keys) {
		return model.User{}, fmt.Errorf("index out of range")
	}

	return it.userInMemory.users[it.keys[it.pos]], nil
}

package service

import "github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/model"

//creational pattern: abstract factory
type FactoryInMemory struct {
}

func (factory *FactoryInMemory) CreateUserService() (User, error) {
	return &UserInMemory{users: map[string]model.User{}}, nil
}

package service

import (
	"testing"
)

func TestFactoryInMemory_CreateUserService(t *testing.T) {

	factory := Factory(&FactoryInMemory{})

	_, err := factory.CreateUserService()

	if err != nil {
		t.Fatal("error not expected", err)
	}
}

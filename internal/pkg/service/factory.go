package service

//creational pattern: abstract factory
type Factory interface {
	CreateUserService() (User, error)
}

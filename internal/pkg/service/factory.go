package service

// abstract factory for creating services
type Factory interface {
	CreateUserService() (User, error)
}

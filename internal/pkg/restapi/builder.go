package restapi

import "github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/service"

//Creational pattern: builder
type Builder interface {
	Reset()
	SetUserService(userService service.User)
}

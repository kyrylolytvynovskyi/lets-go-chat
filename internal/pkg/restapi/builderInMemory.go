package restapi

import "github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/service"

//Creational pattern: builder
type BuilderInMemory struct {
	server *server
}

func (builder *BuilderInMemory) Reset() {
	builder.server = &server{}
}

func (builder *BuilderInMemory) SetUserService(userService service.User) {
	builder.server.userService = userService
}

func (builder *BuilderInMemory) GetProduct() *server {
	srv := builder.server
	builder.Reset()
	return srv
}

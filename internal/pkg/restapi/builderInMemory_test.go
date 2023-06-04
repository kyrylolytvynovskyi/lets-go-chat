package restapi

import (
	"testing"

	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/service"
)

func TestBuilderInMemory_GetProduct(t *testing.T) {
	builder := &BuilderInMemory{}
	//BuildRestServer(builder)
	builder.Reset()
	factory := service.Factory(&service.FactoryInMemory{})
	userService, _ := factory.CreateUserService()
	builder.SetUserService(userService)
	srv := builder.GetProduct()

	if srv == nil {
		t.Fatal("srv must not be null")
	}
}

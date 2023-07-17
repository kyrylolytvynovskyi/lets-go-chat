//go:build wireinject
// +build wireinject

package restapi2

import (
	"github.com/google/wire"
	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/service"
)

func InitializeServer(wsAddr string) *server {
	wire.Build(newServer, service.NewUserInMemory)

	return &server{}
}

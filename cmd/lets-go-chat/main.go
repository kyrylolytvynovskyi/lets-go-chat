package main

import (
	"log"

	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/restapi"
)

func main() {

	log.Fatal(restapi.Run("0.0.0.0:8080"))
}

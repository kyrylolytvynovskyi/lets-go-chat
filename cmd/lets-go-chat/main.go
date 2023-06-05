package main

import (
	"log"

	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/restapi2"
)

func main() {

	log.Fatal(restapi2.Run("0.0.0.0:8080"))
}

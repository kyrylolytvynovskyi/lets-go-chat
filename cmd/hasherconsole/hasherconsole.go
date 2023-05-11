package main

import (
	"fmt"

	"github.com/kyrylolytvynovskyi/lets-go-chat/pkg/hasher"
)

func main() {
	var password string
	fmt.Print("input password:")
	fmt.Scanln(&password)

	hash, err := hasher.HashPassword(password) // ignore error for the sake of simplicity

	if err != nil {
		fmt.Println("Error during password hashing", err)
	}

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := hasher.CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)
}

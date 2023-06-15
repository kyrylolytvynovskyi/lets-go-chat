package hasher

import (
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	myPassword := "myPassword"
	hash, err := HashPassword(myPassword)

	if err != nil {
		t.Error("HashPassword returned an error", err)
		t.FailNow()
	}

	got := CheckPasswordHash(myPassword, hash)
	if !got {
		t.Error("hash and password not equivalent")
	}

	got = CheckPasswordHash(myPassword, "wrong_hash")
	if got {
		t.Error("wrong has and password are equivalent")
	}
}

func ExampleHashPassword() {
	myPassword := "myPassword"
	hash, _ := HashPassword(myPassword)
	isEqual := CheckPasswordHash(myPassword, hash)
	fmt.Println(isEqual)
	// Output:
	// true
}

func ExampleCheckPasswordHash() {
	myPassword := "myPassword"
	fmt.Println(CheckPasswordHash(myPassword, "$2a$14$fo41Lff9uGg3Bmm6OZ1g9uZGsZALyX.8GGwG/Gd0zmngNBM.4RKbG"))
	fmt.Println(CheckPasswordHash(myPassword, "wrong_hash"))
	// Output:
	// true
	// false
}

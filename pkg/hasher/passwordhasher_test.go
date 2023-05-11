package hasher

import "testing"

func TestHashPassword(t *testing.T) {

	myPassword := "myPassword"
	hash, err := HashPassword(myPassword)

	if err != nil {
		t.Error("HashPassword returned an error", err)
		t.FailNow()
	}

	got := CheckPasswordHash(myPassword, hash)
	if !got {
		t.Error("has and password not equivalent")
	}
}

package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChat_GetNewUrl(t *testing.T) {

	chat := NewChat("localhost:8080")

	wsUrl := chat.GetNewUrl("user")

	assert.Contains(t, wsUrl, "localhost:8080")
	tokens := strings.Split(wsUrl, "?token=")
	assert.Equal(t, 2, len(tokens))
}

func TestChat_LoginToken(t *testing.T) {
	//set test data
	chat := NewChat("localhost:8080")
	chat.tokens["valid_token"] = "user"

	//invalid token test
	err := chat.LoginToken("invalid_token")
	assert.ErrorContains(t, err, "invalid token")

	//valid token test
	err = chat.LoginToken("valid_token")
	assert.NoError(t, err)

	//valid token test - relogin
	err = chat.LoginToken("valid_token")
	assert.ErrorContains(t, err, "invalid token")
}

func TestChat_LogoutToken(t *testing.T) {
	//set test data
	chat := NewChat("localhost:8080")
	chat.activeUsers["valid_token"] = "user"

	//logout invalid token
	err := chat.LogoutToken("invalid_token")
	assert.ErrorContains(t, err, "invalid token")

	//logout valid token
	err = chat.LogoutToken("valid_token")
	assert.NoError(t, err)
}

func TestChat_GetActiveUsers(t *testing.T) {
	//set test data
	chat := NewChat("localhost:8080")

	users := chat.GetActiveUsers()
	assert.Len(t, users, 0)

	chat.activeUsers["token0"] = "user0"
	users = chat.GetActiveUsers()
	assert.Len(t, users, 1)
	assert.Equal(t, "user0", users[0])

	chat.activeUsers["token1"] = "user1"
	users = chat.GetActiveUsers()
	assert.Len(t, users, 2)
}

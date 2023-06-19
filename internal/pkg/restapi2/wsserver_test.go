package restapi2

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func Test_server_wsEndpoint(t *testing.T) {

	srv := newServer("localhost:8080", nil)
	tokenUrl := srv.chatService.GetNewUrl("user")
	tokens := strings.Split(tokenUrl, "?token=")
	assert.Len(t, tokens, 2)
	token := tokens[1]

	testServer := httptest.NewServer(http.HandlerFunc(srv.wsEndpoint))
	defer testServer.Close()

	// Invalid token test
	wsUrl := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/ws?token=" + "invalid_token"
	_, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	assert.Error(t, err)

	// Valid token test
	wsUrl = "ws" + strings.TrimPrefix(testServer.URL, "http") + "/ws?token=" + token
	ws, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

}

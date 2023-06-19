package restapi2

import (
	"encoding/json"
	"io/ioutil"
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
	if !assert.NoError(t, err) {
		return
	}
	defer ws.Close()

	err = ws.WriteMessage(websocket.TextMessage, []byte("hello websocket"))
	if !assert.NoError(t, err) {
		return
	}

	_, p, err := ws.ReadMessage()
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "hello websocket", string(p))
}

func Test_server_getActiveUsers(t *testing.T) {

	srv := newServer("localhost:8080", nil)
	srv.routes()

	req := httptest.NewRequest(http.MethodGet, "/v1/users/active", nil)
	resp := httptest.NewRecorder()

	srv.router.ServeHTTP(resp, req)
	result := resp.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusOK, resp.Code)

	data, _ := ioutil.ReadAll(result.Body)

	var activeUsers []string
	json.Unmarshal(data, &activeUsers)

	assert.Len(t, activeUsers, 0)
}

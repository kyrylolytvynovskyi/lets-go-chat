package restapi2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mw "github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/middleware"
	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/model"
	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/service"
	servicemocks "github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/service/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_newInMemoryServer(t *testing.T) {

	server := newInMemoryServer("localhost:8080")
	assert.NotNil(t, server.chatService)
	assert.NotNil(t, server.userService)
	assert.NotNil(t, server.router)
	assert.NotNil(t, server.upgrader)
}

func newMockServer(srv service.User) *server {
	wsAddr := "localhost:8080"

	return newServer(wsAddr, srv)
}

func Test_server_getIndex(t *testing.T) {

	srv := newMockServer(&servicemocks.User{})
	srv.routes()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	srv.router.ServeHTTP(resp, req)
	result := resp.Result()
	defer result.Body.Close()

	if resp.Code != http.StatusOK {
		t.Errorf("HTTP OK expected, but got %v", resp.Code)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	wantGetIndexResponse := "Hello from http server"
	if string(data) != wantGetIndexResponse {
		t.Errorf("expected %v got %v", wantGetIndexResponse, string(data))
	}
}

func testPostUserSuccess(t *testing.T, userServer *servicemocks.User, mux *http.ServeMux) {
	createUserRequest := model.CreateUserRequest{UserName: "new_username", Password: "new_password"}
	userServer.On("CreateUser", createUserRequest).Return(model.CreateUserResponse{Id: "new_id", UserName: "new_username"}, nil)

	jsonReq, _ := json.Marshal(createUserRequest)

	req := httptest.NewRequest(http.MethodPost, "/v1/user", bytes.NewBuffer(jsonReq))
	resp := httptest.NewRecorder()

	mux.ServeHTTP(resp, req)
	result := resp.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusOK, resp.Code, "Wrong http status code")

	data, _ := ioutil.ReadAll(resp.Body)

	createUserResponse := model.CreateUserResponse{}
	json.Unmarshal(data, &createUserResponse)

	assert.Equal(t, "new_id", createUserResponse.Id, "Wrong user id")
	assert.Equal(t, "new_username", createUserResponse.UserName, "Wrong user name")
}

func testPostUserServiceError(t *testing.T, userServer *servicemocks.User, mux *http.ServeMux) {
	//test service error
	createUserRequest := model.CreateUserRequest{UserName: "new_username_service_error", Password: "new_password"}
	userServer.On("CreateUser", createUserRequest).Return(model.CreateUserResponse{}, fmt.Errorf("service error"))

	jsonReq, _ := json.Marshal(createUserRequest)

	req := httptest.NewRequest(http.MethodPost, "/v1/user", bytes.NewBuffer(jsonReq))
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)
	result := resp.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.Code, "Wrong http status code")

	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "service error\n", string(data), "Wrong response text")
}

func testPostUserValidationError(t *testing.T, mux *http.ServeMux) {
	createUserRequest := model.CreateUserRequest{UserName: "usr", Password: "new_password"}
	jsonReq, _ := json.Marshal(createUserRequest)

	req := httptest.NewRequest(http.MethodPost, "/v1/user", bytes.NewBuffer(jsonReq))

	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)
	result := resp.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.Code, "Wrong http status code")
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "userName minLength is 4, actual length is 3\n", string(data), "Wrong response text")
}

func testPostUserParsingError(t *testing.T, mux *http.ServeMux) {

	req := httptest.NewRequest(http.MethodPost, "/v1/user", bytes.NewBuffer([]byte("invalid json")))

	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)
	result := resp.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.Code, "Wrong http status code")
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "invalid character 'i' looking for beginning of value\n", string(data), "Wrong response text")
}

func Test_server_postUser(t *testing.T) {
	userServer := &servicemocks.User{}
	srv := newMockServer(userServer)
	srv.routes()

	testPostUserSuccess(t, userServer, srv.router)
	testPostUserServiceError(t, userServer, srv.router)
	testPostUserValidationError(t, srv.router)
	testPostUserParsingError(t, srv.router)
}

func testPostUserLoginSuccess(t *testing.T, userServer *servicemocks.User, mux *http.ServeMux) {
	reqModel := model.LoginUserRequest{UserName: "user", Password: "password"}
	userServer.On("LoginUser", reqModel).Return(model.LoginUserResponse{"ws://localhost:8080"}, nil)
	jsonReq, _ := json.Marshal(reqModel)

	req := httptest.NewRequest(http.MethodPost, "/v1/user/login", bytes.NewBuffer(jsonReq))

	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)
	result := resp.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusOK, resp.Code, "Wrong http status code")

	data, _ := ioutil.ReadAll(resp.Body)
	var respModel model.LoginUserResponse
	json.Unmarshal(data, &respModel)

	assert.Contains(t, respModel.Url, "ws://localhost:8080", "Wrong websocket url")
}

func testPostUserLoginServiceError(t *testing.T, userServer *servicemocks.User, mux *http.ServeMux) {
	reqModel := model.LoginUserRequest{UserName: "user_error", Password: "password"}
	userServer.On("LoginUser", reqModel).Return(model.LoginUserResponse{}, fmt.Errorf("service error"))
	jsonReq, _ := json.Marshal(reqModel)

	req := httptest.NewRequest(http.MethodPost, "/v1/user/login", bytes.NewBuffer(jsonReq))

	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)
	result := resp.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.Code, "Wrong http status code")

	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "service error\n", string(data), "Wrong response body")

}

func testPostUserLoginParsingError(t *testing.T, mux *http.ServeMux) {

	req := httptest.NewRequest(http.MethodPost, "/v1/user/login", bytes.NewBuffer([]byte("invalid json")))

	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)
	result := resp.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.Code, "Wrong http status code")
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "invalid character 'i' looking for beginning of value\n", string(data), "Wrong response text")
}

func Test_server_postUserLogin(t *testing.T) {
	userServer := &servicemocks.User{}
	srv := newMockServer(userServer)
	srv.routes()

	testPostUserLoginSuccess(t, userServer, srv.router)
	testPostUserLoginServiceError(t, userServer, srv.router)
	testPostUserLoginParsingError(t, srv.router)

}

func Test_server_getError(t *testing.T) {
	srv := newMockServer(&servicemocks.User{})
	srv.routes()

	req := httptest.NewRequest(http.MethodGet, "/error", nil)

	resp := httptest.NewRecorder()
	srv.router.ServeHTTP(resp, req)
	result := resp.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	data, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "test error response")
}

func Test_server_getStringPanic(t *testing.T) {
	srv := newMockServer(&servicemocks.User{})
	srv.routes()

	req := httptest.NewRequest(http.MethodGet, "/panic/string", nil)

	resp := httptest.NewRecorder()

	panicFunc := func() { srv.router.ServeHTTP(resp, req) }

	assert.PanicsWithValue(t, "panic in stringPanic", panicFunc)
}

func Test_server_getStructPanic(t *testing.T) {
	srv := newMockServer(&servicemocks.User{})
	srv.routes()

	req := httptest.NewRequest(http.MethodGet, "/panic/struct", nil)

	resp := httptest.NewRecorder()

	panicFunc := func() { srv.router.ServeHTTP(resp, req) }

	assert.PanicsWithValue(t, mw.PanicStruct{Code: 404, Msg: "panic message"}, panicFunc)
}

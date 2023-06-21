package restapi2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/model"
	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/service"

	mw "github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/middleware"
)

type server struct {
	userService service.User
	chatService *service.Chat
	upgrader    websocket.Upgrader
	router      *http.ServeMux
}

func newServer(wsAddr string, srv service.User) *server {
	//factory := service.Factory(&service.FactoryInMemory{})
	userService := srv
	chatService := service.NewChat("ws://" + wsAddr + "/ws")
	upgrader := websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	return &server{userService: userService, chatService: chatService, upgrader: upgrader, router: http.NewServeMux()}
}

func newInMemoryServer(wsAddr string) *server {
	return newServer(wsAddr, service.NewUserInMemory())
}

func (srv *server) getIndex(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello from http server")
}

func (srv *server) getError(w http.ResponseWriter, r *http.Request) {

	http.Error(w, "test error response", http.StatusInternalServerError)
}

func (srv *server) getStringPanic(w http.ResponseWriter, r *http.Request) {

	panic("panic in stringPanic")
}

func (srv *server) getStructPanic(w http.ResponseWriter, r *http.Request) {

	panic(mw.PanicStruct{Code: 404, Msg: "panic message"})
}

func (srv *server) postUser(w http.ResponseWriter, r *http.Request) {

	var createUserRequest model.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := createUserRequest.Validate(); err != nil {
		log.Printf("Validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createUserResponse, err := srv.userService.CreateUser(createUserRequest)
	if err != nil {
		log.Printf("CreateUser error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(createUserResponse)
}

func (srv *server) postUserLogin(w http.ResponseWriter, r *http.Request) {

	var loginUserRequest model.LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&loginUserRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loginUserResponse, err := srv.userService.LoginUser(loginUserRequest)
	if err != nil {
		log.Printf("LoginUser error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loginUserResponse.Url = srv.chatService.GetNewUrl(loginUserRequest.UserName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginUserResponse)
}

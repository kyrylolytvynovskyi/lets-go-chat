package restapi2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/model"
	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/service"
)

type server struct {
	userService service.User
}

func NewServer(factory service.Factory) *server {
	userService, _ := factory.CreateUserService()

	return &server{userService: userService}
}

func (srv *server) getIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Get is supported", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Hello from http server")
}

func (srv *server) postUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST is supported", http.StatusMethodNotAllowed)
		return
	}

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
	if r.Method != "POST" {
		http.Error(w, "POST is supported", http.StatusMethodNotAllowed)
		return
	}

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginUserResponse)
}

func Run(addr string) error {

	srv := NewServer(&service.FactoryInMemory{})

	http.HandleFunc("/", srv.getIndex)
	http.HandleFunc("/v1/user", srv.postUser)
	http.HandleFunc("/v1/user/login", srv.postUserLogin)

	return http.ListenAndServe(addr, nil)
}

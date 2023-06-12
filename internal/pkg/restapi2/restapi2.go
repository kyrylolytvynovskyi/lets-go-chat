package restapi2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/model"
	"github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/service"

	"github.com/gorilla/handlers"
	mw "github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/middleware"
)

type server struct {
	userService service.User
	router      *http.ServeMux
}

func NewServer(factory service.Factory) *server {
	userService, _ := factory.CreateUserService()

	return &server{userService: userService, router: http.NewServeMux()}
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginUserResponse)
}

func setupRouter() http.Handler {
	srv := NewServer(&service.FactoryInMemory{})

	srv.routes()

	router := mw.ErrorLoggingHandler(os.Stdout)(
		handlers.LoggingHandler(os.Stdout,
			mw.RecoverPanic(srv.router)))

	return router
}

func Run(addr string) error {

	router := setupRouter()

	return http.ListenAndServe(addr, router )
}

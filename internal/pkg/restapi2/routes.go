package restapi2

import (
	"net/http"

	"github.com/gorilla/handlers"
	mw "github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/middleware"
)

func (srv *server) routes() {

	srv.router.Handle("/",
		handlers.MethodHandler{http.MethodGet: http.HandlerFunc(srv.getIndex)})

	srv.router.Handle("/v1/user",
		handlers.MethodHandler{http.MethodPost: http.HandlerFunc(srv.postUser)})

	srv.router.HandleFunc("/v1/user/login",
		mw.EnforceMethod(http.MethodPost,
			srv.postUserLogin))

	srv.router.HandleFunc("/panic/string",
		mw.EnforceMethod(http.MethodGet,
			srv.getStringPanic))

	srv.router.HandleFunc("/panic/struct",
		mw.EnforceMethod(http.MethodGet,
			mw.RecoverPanic(
				srv.getStructPanic)))

}

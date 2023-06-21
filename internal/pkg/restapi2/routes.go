package restapi2

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	mw "github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/middleware"
)

func (srv *server) routes() {

	srv.router.Handle("/",
		mw.EnforceMethod(http.MethodGet,
			srv.getIndex))

	srv.router.Handle("/v1/user",
		mw.EnforceMethod(http.MethodPost,
			srv.postUser))

	srv.router.HandleFunc("/v1/user/login",
		mw.EnforceMethod(http.MethodPost,
			srv.postUserLogin))

	srv.router.Handle("/error",
		mw.EnforceMethod(http.MethodGet,
			srv.getError))

	srv.router.HandleFunc("/panic/string",
		mw.EnforceMethod(http.MethodGet,
			srv.getStringPanic))

	srv.router.HandleFunc("/panic/struct",
		mw.EnforceMethod(http.MethodGet,
			srv.getStructPanic))

	srv.router.HandleFunc("/ws", srv.wsEndpoint)

	srv.router.HandleFunc("/v1/users/active",
		mw.EnforceMethod(http.MethodGet,
			srv.getActiveUsers))

}

func setupRouter(wsAddr string) http.Handler {
	srv := newInMemoryServer(wsAddr)

	srv.routes()

	/*	router := mw.ErrorLoggingHandler(os.Stdout)(
		handlers.LoggingHandler(os.Stdout,
			mw.RecoverPanic(srv.router)))
	*/

	router := handlers.LoggingHandler(os.Stdout,
		mw.RecoverPanic(srv.router))

	return router
}

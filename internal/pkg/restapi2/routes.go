package restapi2

import (
	"net/http"

	mw "github.com/kyrylolytvynovskyi/lets-go-chat/internal/pkg/middleware"
)

func (srv *server) routes() {

	srv.router.HandleFunc("/",
		mw.EnforceMethod(http.MethodGet,
			srv.getIndex))

	srv.router.HandleFunc("/v1/user",
		mw.EnforceMethod(http.MethodPost,
			srv.postUser))

	srv.router.HandleFunc("/v1/user/login",
		mw.EnforceMethod(http.MethodPost,
			srv.postUserLogin))
}

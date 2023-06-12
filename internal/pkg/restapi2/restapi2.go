package restapi2

import "net/http"

func Run(addr, wsAddr string) error {

	router := setupRouter(wsAddr)

	return http.ListenAndServe(addr, router)
}

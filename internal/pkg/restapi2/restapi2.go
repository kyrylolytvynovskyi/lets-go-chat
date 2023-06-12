package restapi2

import "net/http"

func Run(addr string) error {

	router := setupRouter()

	return http.ListenAndServe(addr, router)
}

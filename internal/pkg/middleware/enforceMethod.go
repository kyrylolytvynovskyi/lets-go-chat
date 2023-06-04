package middleware

import (
	"net/http"
)

func EnforceMethod(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != method {
			http.Error(w, "only method"+method+"is supported ", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"net/http"
)

func EnforceMethod(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != method {
			http.Error(w, "only method "+method+" is supported ", http.StatusMethodNotAllowed)
			return
		}

		next(w, r)
	}
}

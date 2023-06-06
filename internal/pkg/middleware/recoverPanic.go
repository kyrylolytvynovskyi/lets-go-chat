package middleware

import (
	"log"
	"net/http"
)

type PanicStruct struct {
	Code int
	Msg  string
}

func RecoverPanic(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic call occured: %v", r)
				switch v := r.(type) {
				case string:
					http.Error(w, v, http.StatusInternalServerError)
				case PanicStruct:
					http.Error(w, v.Msg, v.Code)
				default:
					http.Error(w, "unknown panic", http.StatusInternalServerError)
				}
			}
		}()

		next(w, r)
	}
}

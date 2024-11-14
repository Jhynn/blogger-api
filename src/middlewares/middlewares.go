package middlewares

import (
	"blogger/src/authentication"
	"blogger/src/responses"
	"log"
	"net/http"
)

// Authenticated verifies if the user is authenticated.
func Authenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := authentication.TokenValidation(r); err != nil {
			responses.Error(w, http.StatusUnauthorized, err)
			return
		}

		next(w, r)
	}
}

// Logger prints out the requested method, URI and host in the terminal.
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%6s %s %s", r.Method, r.Host, r.RequestURI)

		next(w, r)
	}
}

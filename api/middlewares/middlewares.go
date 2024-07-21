package middlewares

import (
	"api-center/api/auth"
	"api-center/api/responses"
	"errors"
	"net/http"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "appllication/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(nex http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		nex(w, r)
	}
}

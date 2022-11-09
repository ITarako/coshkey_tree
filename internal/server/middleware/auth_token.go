package middleware

import (
	"errors"
	"net/http"

	"coshkey_tree/internal/server/helper"
)

func AuthToken(next http.Handler, token string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestToken := r.Header.Get("X-Api-Token")
		if requestToken != token {
			serverhelper.WriteError(w, nil, http.StatusForbidden, errors.New("wrong token"), "middleware.AuthToken")
			return
		}

		next.ServeHTTP(w, r)
	})
}

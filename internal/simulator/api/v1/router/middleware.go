package router

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (h *Handler) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHead := r.Header.Get("Authorization")
		authSlice := strings.Split(authHead, " ")
		if !(r.URL.Path == "/sim/auth/sign-up" || r.URL.Path == "/sim/auth/sign-in") {
			if len(authSlice) < 2 || authSlice[0] != "Bearer" {
				errorJSON(w, errors.New("invalid authorisation header"), http.StatusBadRequest)
				return
			}
			token := authSlice[1]
			if err := h.services.ValidateSession(token); err != nil {
				errorJSON(w, errors.New("invalid authorisation token"), http.StatusBadRequest)
				return
			}
		} else {
			if len(authSlice) > 1 || authSlice[0] != "" {
				fmt.Println(authSlice)
				errorJSON(w, errors.New("you already logged in"), http.StatusBadRequest)
				return
			}
		}
		handler.ServeHTTP(w, r)
		return
	})
}

func UpdateSession() {

}

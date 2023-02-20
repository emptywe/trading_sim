package router

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func (h *Handler) middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHead := r.Header.Get("Authorization")
		authSlice := strings.Split(authHead, " ")
		if r.URL.Path == "/sim/auth/refresh" {
			handler.ServeHTTP(w, r)
			return
		}
		if !(r.URL.Path == "/sim/auth/sign-up" || r.URL.Path == "/sim/auth/sign-in") {
			if len(authSlice) < 2 || authSlice[0] != "Bearer" {
				zap.S().Error("invalid authorisation header")
				errorJSON(w, errors.New("invalid authorisation header"), http.StatusBadRequest)
				return
			}
			token := authSlice[1]
			if err := h.services.ValidateSession(token); err != nil {
				zap.S().Errorf("invalid authorisation token: %v", err)
				errorJSON(w, errors.New("invalid authorisation token"), http.StatusBadRequest)
				return
			}
		} else {
			if len(authSlice) > 1 || authSlice[0] != "" {
				zap.S().Error("user already logged in")
				errorJSON(w, errors.New("you already logged in"), http.StatusBadRequest)
				return
			}
		}
		handler.ServeHTTP(w, r)
		return
	})
}

func (h *Handler) refreshSession(w http.ResponseWriter, r *http.Request) {

	headInfo := r.Header.Get("Authorization")
	if headInfo == "" {
		zap.S().Error("invalid authorisation header")
		errorJSON(w, errors.New("invalid authorisation header"), http.StatusBadRequest)
		return
	}
	tokenStruct := strings.Split(headInfo, " ")
	expiredToken := tokenStruct[1]

	var tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(r.Body).Decode(&tokenReqBody)
	if err != nil {
		zap.S().Errorf("can't parse request token: %v", err)
		errorJSON(w, err, http.StatusBadRequest)
		return
	}

	token, err := h.services.UpdateSession(expiredToken, tokenReqBody.RefreshToken)
	if err != nil {
		zap.S().Errorf("can't update session: %v", err)
		errorJSON(w, err, http.StatusBadRequest)
		return
	}
	_ = json.NewEncoder(w).Encode(token)
}

package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/emptywe/trading_sim/entity"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {

	var user entity.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		zap.S().Errorf("unable to decode json: %v", err)
		errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := user.ValidateUser(); err != nil {
		zap.S().Errorf("can't validate user: %v", err)
		errorJSON(w, err, http.StatusLengthRequired)
		return
	}

	id, err := h.services.Authorization.CreateUser(user)
	if err != nil {
		err = fmt.Errorf("user not created: %v", err)
		zap.S().Error(err)
		errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err = h.services.Basket.CreateStartingBasket(id); err != nil {
		err = fmt.Errorf("basket not created, user dropped: %v", err)
		zap.S().Error(err)
		errorJSON(w, err, http.StatusInternalServerError)
		_ = h.services.Authorization.DeleteUser(id)
		return
	}

	if err = json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	}); err != nil {
		zap.S().Error("can't send signUp success response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {

	var request entity.SignInRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		zap.S().Errorf("can't unmarshal json: %v", err)
		errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.services.ReadUser(request)
	if err != nil {
		zap.S().Errorf("can't read user %v", err)
		errorJSON(w, err, http.StatusBadRequest)
		return
	}

	session, err := h.services.CreateSession(&user)
	if err != nil {
		zap.S().Errorf("can't create session %v", err)
		errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(w).Encode(entity.SignInResponse{
		User:    user,
		Session: session,
	}); err != nil {
		zap.S().Error("can't send signIn success response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {

	authHead := r.Header.Get("Authorization")
	authSlice := strings.Split(authHead, " ")
	token := authSlice[1]
	err := h.services.DeleteSession(token)
	if err != nil {
		zap.S().Errorf("can't delete session: %v", err)
		errorJSON(w, err, http.StatusBadRequest)
		return
	}
	if err = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "successfully logged out",
	}); err != nil {
		zap.S().Error("can't send signUp success response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {

}

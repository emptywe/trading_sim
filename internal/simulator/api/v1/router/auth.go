package router

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"

	"github.com/emptywe/trading_sim/model"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {

	var user model.SignUpRequest
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

	id, err := h.services.Authorization.CreateUser(user.Email, user.UserName, user.Password)
	if err != nil {
		err = fmt.Errorf("user not created: %v", err)
		zap.S().Error(err)
		errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	_, err = h.services.Basket.CreateStartingBasket(id)
	if err != nil {
		err = fmt.Errorf("basket not created, user dropped: %v", err)
		zap.S().Error(err)
		errorJSON(w, err, http.StatusInternalServerError)
		_ = h.services.Authorization.DropUser(id)
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

type signInInput struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {

	var input signInInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		zap.S().Errorf("can't unmarshal json: %v", err)
		errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// TODO: add create session logic

}

func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {

	//co, err := c.Request.Cookie("USession")
	//if err != nil {
	//	c.AbortWithStatus(http.StatusUnauthorized)
	//	return
	//}
	//
	//_, ok := h.services.Authorization.ValidateSession(co)
	//if !ok {
	//	c.AbortWithStatus(http.StatusUnauthorized)
	//	return
	//}
	//
	//_ = h.services.Authorization.ExpireSession(co)
	//co.Path = "/sim"
	//co.MaxAge = -1
	//
	//http.SetCookie(c.Writer, co)

}

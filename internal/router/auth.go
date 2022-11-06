package router

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/emptywe/trading_sim/model"
)

func (h *Handler) signUp(c *gin.Context) {

	_, err := c.Request.Cookie("USession")
	if err == nil {
		newErrorResponse(c, http.StatusConflict, "You logged in")
		return
	}

	var req model.SignUpRequest

	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if len(req.UserName) < 2 {
		err := errors.New("Username must be at least two sybols")
		c.JSON(http.StatusOK, map[string]interface{}{
			"error": err,
		})
		return
	}

	if len(req.Password) < 4 {
		err := errors.New("Short password")
		c.JSON(http.StatusOK, map[string]interface{}{
			"error": err,
		})
		return
	}

	id, err := h.services.Authorization.CreateUser(req.Email, req.UserName, req.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	_, err = h.services.Basket.CreateStartingBasket(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		_ = h.services.Authorization.DropUser(id)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {

	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err := c.Request.Cookie("USession")
	if err == nil {
		newErrorResponse(c, http.StatusConflict, "You already logged in")
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.UserName, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	co := &http.Cookie{
		Name:   "USession",
		Value:  token,
		Path:   "/sim",
		MaxAge: 7200,
	}

	id, suid, err := h.services.CreateSession(co)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(c.Writer, co)

	c.JSON(http.StatusOK, map[string]interface{}{
		"token":   token,
		"id":      id,
		"session": suid,
	})

}

func (h *Handler) logOut(c *gin.Context) {

	co, err := c.Request.Cookie("USession")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	_, ok := h.services.Authorization.ValidateSession(co)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	_ = h.services.Authorization.ExpireSession(co)
	co.Path = "/sim"
	co.MaxAge = -1

	http.SetCookie(c.Writer, co)

}

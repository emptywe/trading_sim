package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	userCtx = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	co, err := c.Request.Cookie("USession")
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	_, ok := h.services.Authorization.ValidateSession(co)
	if !ok {
		co.MaxAge = -1
		http.SetCookie(c.Writer, co)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//c.JSON(http.StatusOK, map[string]interface{}{
	//	"id" :id,
	//	"ok": ok,
	//})

	co.Path = "/sim"
	co.MaxAge = 7200

	http.SetCookie(c.Writer, co)

}

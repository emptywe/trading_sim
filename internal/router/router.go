package router

import (
	"github.com/emptywe/trading_sim/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     *service.Service
	currencyList []string
}

func NewHandler(services *service.Service, currencyList []string) *Handler {
	return &Handler{services: services, currencyList: currencyList}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.HandleMethodNotAllowed = true
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	sim := router.Group("/sim")
	{
		auth := sim.Group("/auth")
		{
			auth.POST("/sign-up", h.signUp)
			auth.POST("/sign-in", h.signIn)
			auth.POST("/logout", h.logOut)
		}

		api := sim.Group("/api", h.userIdentity)

		{
			basket := api.Group("/basket")
			{
				basket.GET("/prices", h.prices)
				basket.POST("/swap", h.swap)
				basket.GET("/balance", h.balance)
				basket.GET("/top", h.topUsers)
			}
		}
	}
	return router
}

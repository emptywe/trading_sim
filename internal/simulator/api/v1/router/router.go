package router

import (
	"github.com/emptywe/trading_sim/internal/simulator/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	services     *service.Service
	currencyList []string
}

func NewHandler(services *service.Service, currencyList []string) *Handler {
	return &Handler{services: services, currencyList: currencyList}
}

func (h *Handler) InitRoutes() http.Handler {
	router := mux.NewRouter()
	router.Use(h.Middleware)
	router.HandleFunc("/sim/auth/sign-up", h.signUp).Methods("POST").Name("SignUp")
	router.HandleFunc("/sim/auth/sign-in", h.signIn).Methods("POST").Name("SignIp")
	router.HandleFunc("/sim/auth/update-user", h.updateUser).Methods("POST").Name("Update")
	router.HandleFunc("/sim/auth/delete-user", nil).Methods("DELETE").Name("Delete")
	router.HandleFunc("/sim/auth/logout", h.logOut).Methods("POST").Name("Logout")
	router.HandleFunc("/sim/api/basket/prices", h.prices).Methods("GET").Name("Prices")
	router.HandleFunc("/sim/api/basket/swap", h.swap).Methods("POST").Name("Swap")
	router.HandleFunc("/sim/api/basket/balance", h.balance).Methods("GET").Name("Balance")
	router.HandleFunc("/sim/api/basket/top", h.topUsers).Methods("GET").Name("TopUsers")
	if err := router.Walk(DebugRouter); err != nil {
		zap.S().Error(err)
	}
	//handler := cors.AllowAll().Handler(router)
	zap.S().Info("Router initialised")
	return handlers.CORS()(router)
}

func DebugRouter(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	path, _ := route.GetPathTemplate()
	if len([]rune(path)) < 1 {
		path = "undefined router path"
	}
	methods, _ := route.GetMethods()
	if len(methods) < 1 {
		route.Methods("GET")
		methods = append(methods, "undefined methods, default GET")
	}
	zap.S().Debugf("route: %s\t\tpath: %s %s", route.GetName(), path, methods)
	return nil
}

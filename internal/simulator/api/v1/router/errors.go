package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func errorJSON(w http.ResponseWriter, err error, statusCode int) {
	if err := json.NewEncoder(w).Encode(map[string]string{
		"code": fmt.Sprintf("%d", statusCode), "error": err.Error(),
	}); err != nil {
		zap.S().Error("can't send error response " + err.Error())
	}
	w.WriteHeader(statusCode)
}

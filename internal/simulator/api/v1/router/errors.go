package router

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

func errorJSON(w http.ResponseWriter, err error, statusCode int) {
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		fmt.Sprintf("%d", statusCode): err,
	}); err != nil {
		zap.S().Error("can't send error response " + err.Error())
	}
	w.WriteHeader(statusCode)
}

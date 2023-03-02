package router

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) prices(w http.ResponseWriter, r *http.Request) {

	carr, err := h.services.Info.GetAllCurrenciesUSD()
	if err != nil {
		zap.S().Errorf("can't get all currencies usd: %v", err)
		errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(map[string]interface{}{
		"currencies": carr,
	}); err != nil {
		zap.S().Error("can't send prices success response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *Handler) topUsers(w http.ResponseWriter, r *http.Request) {

	tu, err := h.services.Info.GetTopUsers()
	if err != nil {
		errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(map[string]interface{}{
		"top": tu,
	}); err != nil {
		zap.S().Error("can't send topUsers success response ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

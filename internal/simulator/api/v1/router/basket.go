package router

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/emptywe/trading_sim/entity"
)

func (h *Handler) swap(w http.ResponseWriter, r *http.Request) {

	var input entity.Transaction
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		zap.S().Errorf("can't decode transaction json: %v", err)
		errorJSON(w, err, http.StatusBadRequest)
		return
	}

	uid := r.Context().Value("uid").(int)
	err := h.services.ServeTrade(input, uid)
	if err != nil {
		zap.S().Errorf("can't sevre trade: %v", err)
		errorJSON(w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) balance(w http.ResponseWriter, r *http.Request) {

	uid := r.Context().Value("uid").(int)
	baskets, err := h.services.Basket.GetAllBaskets(uid)
	if err != nil {
		zap.S().Errorf("can't' get all baskets: %v", err)
		errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var balance float64
	for i, _ := range baskets {
		balance += baskets[i].USDAmount
	}

	if err = json.NewEncoder(w).Encode(map[string]interface{}{
		"balance": balance,
		"baskets": baskets,
	}); err != nil {
		zap.S().Error("can't send balance success response ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

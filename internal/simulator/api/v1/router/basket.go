package router

import (
	"encoding/json"
	"github.com/emptywe/trading_sim/entity"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) swap(w http.ResponseWriter, r *http.Request) {

	var input entity.Transaction
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		zap.S().Errorf("can't decode transaction json: %v", err)
		errorJSON(w, err, http.StatusBadRequest)
		return
	}
	// TODO: move to another function
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

	// TODO: create new user balance logic
	bb, err := h.services.Basket.GetAllBaskets(1)
	if err != nil {
		zap.S().Errorf("can't' get all baskets: %v", err)
		errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var b float64
	for i, _ := range bb {
		b += bb[i].USDAmount
	}

	if err = json.NewEncoder(w).Encode(map[string]interface{}{
		"balance": b,
		"baskets": bb,
	}); err != nil {
		zap.S().Error("can't send balance success response ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

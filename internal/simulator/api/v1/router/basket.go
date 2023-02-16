package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Transaction struct {
	C1Name string  `json:"c1_name"`
	C2Name string  `json:"c2_name"`
	Dvalue float64 `json:"dvalue"`
}

func (h *Handler) prices(w http.ResponseWriter, r *http.Request) {

	carr, err := h.services.Basket.GetAllCurrenciesUSD()
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

func (h *Handler) swap(w http.ResponseWriter, r *http.Request) {

	var input Transaction
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		zap.S().Errorf("can't decode transaction json: %v", err)
		errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//id, _, err := h.services.Authorization.ParseToken(co.Value)
	id := 1 // TODO: add users session_cache logic

	// TODO: move to anothre function
	if input.C2Name == "usdt" {
		bcur := fmt.Sprintf(input.C1Name + input.C2Name)
		fmt.Println(bcur)

		bId, err := h.services.Basket.CreateBasketSell(id, bcur, input.Dvalue)
		if err != nil {
			zap.S().Errorf("can't create basket sell: %v", err)
			errorJSON(w, err, http.StatusInternalServerError)
			return
		}
		fmt.Println(err)
		if err = json.NewEncoder(w).Encode(map[string]interface{}{
			"bId": bId,
		}); err != nil {
			zap.S().Errorf("can't send create basket sell success response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		bcur := fmt.Sprintf(input.C2Name + input.C1Name)
		fmt.Println(bcur)

		bId, err := h.services.Basket.CreateBasket(id, input.C1Name, bcur, input.Dvalue)
		if err != nil {
			zap.S().Errorf("can't create basket: %v", err)
			errorJSON(w, err, http.StatusInternalServerError)
			return
		}
		fmt.Println(err)
		if err = json.NewEncoder(w).Encode(map[string]interface{}{
			"bId": bId,
		}); err != nil {
			zap.S().Errorf("can't send create basket sell success response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) balance(w http.ResponseWriter, r *http.Request) {

	//TODO: move to another function
	for _, cur := range h.currencyList {
		err := h.services.Basket.UpdateBasket(cur)
		if err != nil {
			zap.S().Errorf("can't' update basket: %v", err)
		}
	}
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

func (h *Handler) topUsers(w http.ResponseWriter, r *http.Request) {

	_, err := h.services.Basket.UpdateBalance()
	if err != nil {
		zap.S().Errorf("can't update user balance error: %v", err)
	}

	tu, err := h.services.Basket.GetTopUsers()
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

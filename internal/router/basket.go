package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Transaction struct {
	C1Name string  `json:"c1_name"`
	C2Name string  `json:"c2_name"`
	Dvalue float64 `json:"dvalue"`
}

func (h *Handler) prices(c *gin.Context) {

	carr, err := h.services.Basket.GetAllCurrenciesUSD()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"currencies": carr,
	})

}

func (h *Handler) swap(c *gin.Context) {

	co, err := c.Request.Cookie("USession")
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		c.AbortWithStatus(http.StatusTemporaryRedirect)
		return
	}

	var input Transaction
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, _, err := h.services.Authorization.ParseToken(co.Value)

	if input.C2Name == "usdt" {
		bcur := fmt.Sprintf(input.C1Name + input.C2Name)
		fmt.Println(bcur)

		bId, err := h.services.Basket.CreateBasketSell(id, bcur, input.Dvalue)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, map[string]string{
				"message": err.Error(),
			})
		}
		fmt.Println(err)
		c.JSON(http.StatusOK, map[string]interface{}{
			"bId": bId,
		})
	} else {
		bcur := fmt.Sprintf(input.C2Name + input.C1Name)
		fmt.Println(bcur)

		bId, err := h.services.Basket.CreateBasket(id, input.C1Name, bcur, input.Dvalue)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, map[string]string{
				"message": err.Error(),
			})
		}
		fmt.Println(err)
		c.JSON(http.StatusOK, map[string]interface{}{
			"bId": bId,
		})

	}

}

func (h *Handler) balance(c *gin.Context) {

	co, err := c.Request.Cookie("USession")
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		c.AbortWithStatus(http.StatusTemporaryRedirect)
		return
	}

	id, _, err := h.services.Authorization.ParseToken(co.Value)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	for _, cur := range h.currencyList {
		err := h.services.Basket.UpdateBasket(cur)
		if err != nil {
			logrus.Printf("can't' update basket, error: %s", err.Error())
		}
	}

	bb, err := h.services.Basket.GetAllBaskets(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var b float64
	for i, _ := range bb {
		b += bb[i].USDAmount
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"balance": b,
		"baskets": bb,
	})

}

func (h *Handler) topUsers(c *gin.Context) {

	_, err := h.services.Basket.UpdateBalance()
	if err != nil {
		logrus.Errorf("can't update user balance error: %s", err.Error())
	}

	tu, err := h.services.Basket.GetTopUsers()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"top": tu,
	})

}

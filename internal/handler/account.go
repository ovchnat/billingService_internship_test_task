package handler

import (
	"billingService/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) depositMoney(c *gin.Context) {
	var input entity.UpdateBalanceRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	accResponse, err := h.services.BalanceOperations.DepositMoney(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"user-id":             accResponse.UserId,
		"account-id":          accResponse.AccountId,
		"user-balance":        accResponse.Balance,
		"user-pending-amount": accResponse.Pending,
	})
}

func (h *Handler) withdrawMoney(c *gin.Context) {
	var input entity.UpdateBalanceRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	accResponse, err := h.services.BalanceOperations.WithdrawMoney(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"user-id":             accResponse.UserId,
		"user-balance":        accResponse.Balance,
		"user-pending-amount": accResponse.Pending,
	})
}

func (h *Handler) getBalance(c *gin.Context) {
	var input entity.GetBalanceRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	accResponse, err := h.services.BalanceOperations.GetBalance(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"user-balance":        accResponse.Balance,
		"user-pending-amount": accResponse.Pending,
	})
}

func (h *Handler) reserveAmount(c *gin.Context) {
	var input entity.ReserveServiceFeeRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	accResponse, err := h.services.BalanceOperations.ReserveServiceFee(input, c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"user-id":    accResponse.UserId,
		"service-id": accResponse.ServiceId,
		"order-id":   accResponse.OrderId,
		"invoice":    accResponse.Invoice,
		"status":     accResponse.Status,
		"created-at": accResponse.CreatedAt,
	})
}

func (h *Handler) orderConfirm(c *gin.Context) {

}

func (h *Handler) transfer(c *gin.Context) {

}

func (h *Handler) sayHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World!",
	})
}

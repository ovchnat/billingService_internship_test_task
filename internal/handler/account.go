package handler

import (
	"billingService/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) depositMoney(c *gin.Context) {
	var input entity.UpdateBalanceRequest

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	accResponse, err := h.services.BalanceOperations.DepositMoney(input, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"account-id":       accResponse.AccountId,
		"sum-deposited":    accResponse.Sum,
		"operation-status": accResponse.Status,
		"operation-event":  accResponse.EventType,
		"created-at":       accResponse.CreatedAt,
	})
}

func (h *Handler) withdrawMoney(c *gin.Context) {
	var input entity.UpdateBalanceRequest

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	accResponse, err := h.services.BalanceOperations.WithdrawMoney(input, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"account-id":       accResponse.AccountId,
		"sum-withdrawn":    accResponse.Sum,
		"operation-status": accResponse.Status,
		"operation-event":  accResponse.EventType,
		"created-at":       accResponse.CreatedAt,
	})
}

func (h *Handler) getBalance(c *gin.Context) {
	var input entity.GetBalanceRequest

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	accResponse, err := h.services.BalanceOperations.GetBalance(input, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"user-balance":        accResponse.Balance,
		"user-pending-amount": accResponse.Pending,
	})
}

func (h *Handler) reserveServiceFee(c *gin.Context) {
	var input entity.ReserveServiceFeeRequest

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	accResponse, err := h.services.BalanceOperations.ReserveServiceFee(input, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"account-id": accResponse.AccountId,
		"service-id": accResponse.ServiceId,
		"order-id":   accResponse.OrderId,
		"invoice":    accResponse.Invoice,
		"status":     accResponse.Status,
		"created-at": accResponse.CreatedAt,
		"updated-at": accResponse.UpdatedAt,
	})
}

func (h *Handler) approveOrderFee(c *gin.Context) {
	var input entity.StatusServiceFeeRequest

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	accResponse, err := h.services.BalanceOperations.ApproveServiceFee(input, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"account-id": accResponse.AccountId,
		"service-id": accResponse.ServiceId,
		"order-id":   accResponse.OrderId,
		"invoice":    accResponse.Invoice,
		"status":     accResponse.Status,
		"created-at": accResponse.CreatedAt,
		"updated-at": accResponse.UpdatedAt,
	})
}

func (h *Handler) failedServiceFee(c *gin.Context) {
	var input entity.StatusServiceFeeRequest

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	accResponse, err := h.services.BalanceOperations.FailedServiceFee(input, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"account-id": accResponse.AccountId,
		"service-id": accResponse.ServiceId,
		"order-id":   accResponse.OrderId,
		"invoice":    accResponse.Invoice,
		"status":     accResponse.Status,
		"created-at": accResponse.CreatedAt,
		"updated-at": accResponse.UpdatedAt,
	})
}

func (h *Handler) transfer(c *gin.Context) {
	var input entity.TransferRequest

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	accResponse, err := h.services.BalanceOperations.Transfer(input, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"receive-account":  accResponse.AccountIdTo,
		"transfer-account": accResponse.AccountIdFrom,
		"amount":           accResponse.Amount,
		"status":           accResponse.Status,
		"event-type":       accResponse.EventType,
		"created-at":       accResponse.Timecode,
	})
}

func (h *Handler) sayHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World!",
	})
}

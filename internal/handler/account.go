package handler

import (
	"billingService/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (h *Handler) getBalance(c *gin.Context) {
	idStringInput := c.Param("id")
	logrus.Printf("Input read: %v %T", idStringInput, idStringInput)

	idNumberInput, err := strconv.ParseInt(idStringInput, 10, 64)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	input := entity.GetBalanceRequest{UserId: idNumberInput}

	response, newErr := h.services.BalanceOperations.GetBalance(input, c)
	if newErr != nil {
		NewErrorResponse(c, http.StatusInternalServerError, newErr.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"user-balance":        response.Balance,
			"user-pending-amount": response.Pending,
		})
	}
}

func (h *Handler) depositMoney(c *gin.Context) {
	var updateBalanceDepositRequest entity.UpdateBalanceRequest

	if err := c.BindJSON(&updateBalanceDepositRequest); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.BalanceOperations.DepositMoney(updateBalanceDepositRequest, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"account-id":       response.AccountId,
			"sum-deposited":    response.Sum,
			"operation-status": response.Status,
			"operation-event":  response.EventType,
			"created-at":       response.CreatedAt,
		})
	}
}

func (h *Handler) withdrawMoney(c *gin.Context) {
	var updateBalanceWithdrawRequest entity.UpdateBalanceRequest

	if err := c.BindJSON(&updateBalanceWithdrawRequest); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.BalanceOperations.WithdrawMoney(updateBalanceWithdrawRequest, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"account-id":       response.AccountId,
			"sum-withdrawn":    response.Sum,
			"operation-status": response.Status,
			"operation-event":  response.EventType,
			"created-at":       response.CreatedAt,
		})
	}
}

func (h *Handler) reserveServiceFee(c *gin.Context) {
	var reserveServiceFeeRequest entity.ReserveServiceFeeRequest

	if err := c.BindJSON(&reserveServiceFeeRequest); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.BalanceOperations.ReserveServiceFee(reserveServiceFeeRequest, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"account-id": response.AccountId,
			"service-id": response.ServiceId,
			"order-id":   response.OrderId,
			"invoice":    response.Invoice,
			"status":     response.Status,
			"created-at": response.CreatedAt,
			"updated-at": response.UpdatedAt,
		})
	}
}

func (h *Handler) approveOrderFee(c *gin.Context) {
	var statusApproveServiceFeeRequest entity.StatusServiceFeeRequest

	if err := c.BindJSON(&statusApproveServiceFeeRequest); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.BalanceOperations.ApproveServiceFee(statusApproveServiceFeeRequest, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"account-id": response.AccountId,
			"service-id": response.ServiceId,
			"order-id":   response.OrderId,
			"invoice":    response.Invoice,
			"status":     response.Status,
			"created-at": response.CreatedAt,
			"updated-at": response.UpdatedAt,
		})
	}
}

func (h *Handler) failedServiceFee(c *gin.Context) {
	var statusFailedServiceFeeRequest entity.StatusServiceFeeRequest

	if err := c.BindJSON(&statusFailedServiceFeeRequest); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.BalanceOperations.FailedServiceFee(statusFailedServiceFeeRequest, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"account-id": response.AccountId,
			"service-id": response.ServiceId,
			"order-id":   response.OrderId,
			"invoice":    response.Invoice,
			"status":     response.Status,
			"created-at": response.CreatedAt,
			"updated-at": response.UpdatedAt,
		})
	}
}

func (h *Handler) transfer(c *gin.Context) {
	var transferRequest entity.TransferRequest

	if err := c.BindJSON(&transferRequest); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.BalanceOperations.Transfer(transferRequest, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"receive-account":  response.AccountIdTo,
			"transfer-account": response.AccountIdFrom,
			"amount":           response.Amount,
			"status":           response.Status,
			"event-type":       response.EventType,
			"created-at":       response.Timecode,
		})
	}
}

func (h *Handler) sayHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World!",
	})
}

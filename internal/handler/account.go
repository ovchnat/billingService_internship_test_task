package handler

import (
	"billingService/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// @Summary getBalance
// @Tags account
// @Description "Deposit money for a given account"
// @Produce json
// @Param id path integer true "user id"
// @Success 200 {object} entity.GetBalanceResponse
// @Failure 500 {object} errorAcc
// @Router /account/getBalance/{id} [get]
func (h *Handler) getBalance(c *gin.Context) {
	idStringInput := c.Param("id")
	logrus.Printf("Input read: %v %T", idStringInput, idStringInput)

	idNumberInput, err := strconv.ParseInt(idStringInput, 10, 64)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
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

// @Summary depositMoney
// @Tags account
// @Description "deposit money for a given account"
// @Accept json
// @Produce json
// @Param input body entity.UpdateBalanceRequest true "JSON object with user ID and money amount to deposit"
// @Success 200 {object} entity.UpdateBalanceDepositResponse
// @Failure 500 {object} errorAcc
// @Router /account/depositMoney [post]
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

// @Summary withdrawMoney
// @Tags account
// @Description "Withdraw money for a given account"
// @Accept json
// @Produce json
// @Param input body entity.UpdateBalanceRequest true "JSON object with user ID and money amount to withdraw"
// @Success 200 {object} entity.UpdateBalanceWithdrawResponse
// @Failure 500 {object} errorAcc
// @Router /account/withdrawMoney [post]
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

// @Summary reserveServiceFee
// @Tags account
// @Description "Put specified amount of money in reservation for a given account"
// @Accept json
// @Produce json
// @Param input body entity.ReserveServiceFeeRequest true "JSON object with used ID, service ID, order ID and fee amount"
// @Success 200 {object} entity.ReserveServiceFeeResponse
// @Failure 500 {object} errorAcc
// @Router /account/reserveServiceFee [post]
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

// @Summary approveOrderFee
// @Tags account
// @Description "Approve specified reservation"
// @Accept json
// @Produce json
// @Param input body entity.StatusServiceFeeRequest true "JSON object with used ID, service ID, order ID and fee amount"
// @Success 200 {object} entity.StatusServiceFeeResponse
// @Failure 500 {object} errorAcc
// @Router /account/approveOrderFee [post]
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

// @Summary failedServiceFee
// @Tags account
// @Description "Mark reservation as failed and release funds"
// @Accept json
// @Produce json
// @Param input body entity.StatusServiceFeeRequest true "JSON object with used ID, service ID, order ID and fee amount"
// @Success 200 {object} entity.StatusServiceFeeResponse
// @Failure 500 {object} errorAcc
// @Router /account/failedServiceFee [post]
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

// @Summary transfer
// @Tags account
// @Description "Transfer funds from one account to another"
// @Accept json
// @Produce json
// @Param input body entity.TransferRequest true "JSON object with sender ID, receiver ID and money amount"
// @Success 200 {object} entity.TransferResponse
// @Failure 500 {object} errorAcc
// @Router /account/transfer [post]
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

package service

import (
	"billingService/internal/entity"
	"billingService/internal/repository"
	"github.com/gin-gonic/gin"
)

type BalanceOperations interface {
	GetBalance(userId entity.GetBalanceRequest, ctx *gin.Context) (entity.GetBalanceResponse, error)
	DepositMoney(depositReq entity.UpdateBalanceRequest, ctx *gin.Context) (entity.UpdateBalanceResponse, error)
	WithdrawMoney(withdrawReq entity.UpdateBalanceRequest, ctx *gin.Context) (entity.UpdateBalanceResponse, error)
	ReserveServiceFee(reserveSerFeeReq entity.ReserveServiceFeeRequest, ctx *gin.Context) (entity.ReserveServiceFeeResponse, error)
	ApproveServiceFee(approveSerFeeReq entity.StatusServiceFeeRequest, ctx *gin.Context) (entity.StatusServiceFeeResponse, error)
	Transfer(transferReq entity.TransferRequest, ctx *gin.Context) (entity.TransferResponse, error)
	FailedServiceFee(failedServiceFeeReq entity.StatusServiceFeeRequest, ctx *gin.Context) (entity.StatusServiceFeeResponse, error)
}

type ReportOperations interface {
	ServiceMonthlyReport(serviceMonthlyReport entity.ServiceMonthlyReportReq, ctx *gin.Context) (entity.ServiceMonthlyReportResponse, error)
}

type BillingService struct {
	BalanceOperations
	ReportOperations
}

func NewService(repos *repository.BillingRepo) *BillingService {
	return &BillingService{
		BalanceOperations: NewBalanceOperationsService(repos.BalanceOperations),
		ReportOperations:  NewReportOperationsService(repos.ReportOperations),
	}
}

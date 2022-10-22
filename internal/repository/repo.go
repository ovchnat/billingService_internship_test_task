package repository

import (
	"billingService/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

type BillingRepo struct {
	BalanceOperations
	ReportOperations
}

func NewRepo(db *sqlx.DB) *BillingRepo {
	return &BillingRepo{
		BalanceOperations: NewAccPostgres(db),
		ReportOperations:  NewReportPostgres(db),
	}
}

package repository

import (
	"billingService/internal/entity"
	"context"
	"github.com/jmoiron/sqlx"
)

type BalanceOperations interface {
	GetBalance(userId entity.GetBalanceRequest) (entity.GetBalanceResponse, error)
	DepositMoney(depositReq entity.UpdateBalanceRequest) (entity.UpdateBalanceResponse, error)
	WithdrawMoney(withdrawReq entity.UpdateBalanceRequest) (entity.UpdateBalanceResponse, error)
	ReserveServiceFee(reserveSerFeeReq entity.ReserveServiceFeeRequest, ctx context.Context) (entity.ReserveServiceFeeResponse, error)
}

type BillingRepo struct {
	BalanceOperations
}

func NewRepo(db *sqlx.DB) *BillingRepo {
	return &BillingRepo{
		BalanceOperations: NewAccPostgres(db),
	}
}

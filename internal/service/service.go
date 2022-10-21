package service

import (
	"billingService/internal/entity"
	"billingService/internal/repository"
)

type BalanceOperations interface {
	GetBalance(userId entity.GetBalanceRequest) (entity.GetBalanceResponse, error)
	DepositMoney(depositReq entity.UpdateBalanceRequest) (entity.UpdateBalanceResponse, error)
	ReserveServiceFee(reserveSerFeeReq entity.ReserveServiceFeeRequest) (entity.ReserveServiceFeeResponse, error)
}

type BillingService struct {
	BalanceOperations
}

func NewService(repos *repository.BillingRepo) *BillingService {
	return &BillingService{
		BalanceOperations: NewBalanceOperationsService(repos.BalanceOperations),
	}
}

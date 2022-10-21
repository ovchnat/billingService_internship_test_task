package service

import (
	"billingService/internal/entity"
	"billingService/internal/repository"
)

type BalanceOperationsService struct {
	repo repository.BalanceOperations
}

func NewBalanceOperationsService(repo repository.BalanceOperations) *BalanceOperationsService {
	return &BalanceOperationsService{repo: repo}
}

func (s *BalanceOperationsService) GetBalance(userid entity.GetBalanceRequest) (entity.GetBalanceResponse, error) {
	return s.repo.GetBalance(userid)
}

func (s *BalanceOperationsService) DepositMoney(depReq entity.UpdateBalanceRequest) (entity.UpdateBalanceResponse, error) {
	return s.repo.DepositMoney(depReq)
}

func (s *BalanceOperationsService) ReserveServiceFee(reserveReq entity.ReserveServiceFeeRequest) (entity.ReserveServiceFeeResponse, error) {
	return s.repo.ReserveServiceFee(reserveReq)
}

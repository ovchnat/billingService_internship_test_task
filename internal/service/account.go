package service

import (
	"billingService/internal/entity"
	"billingService/internal/repository"
	"github.com/gin-gonic/gin"
)

type BalanceOperationsService struct {
	repo repository.BalanceOperations
}

func NewBalanceOperationsService(repo repository.BalanceOperations) *BalanceOperationsService {
	return &BalanceOperationsService{repo: repo}
}

func (s *BalanceOperationsService) GetBalance(userid entity.GetBalanceRequest, ctx *gin.Context) (entity.GetBalanceResponse, error) {
	return s.repo.GetBalance(userid, ctx)
}

func (s *BalanceOperationsService) DepositMoney(depReq entity.UpdateBalanceRequest, ctx *gin.Context) (entity.UpdateBalanceResponse, error) {
	return s.repo.DepositMoney(depReq, ctx)
}

func (s *BalanceOperationsService) ReserveServiceFee(reserveReq entity.ReserveServiceFeeRequest, ctx *gin.Context) (entity.ReserveServiceFeeResponse, error) {
	return s.repo.ReserveServiceFee(reserveReq, ctx)
}

func (s *BalanceOperationsService) ApproveServiceFee(appSerFeeReq entity.StatusServiceFeeRequest, ctx *gin.Context) (entity.StatusServiceFeeResponse, error) {
	return s.repo.ApproveServiceFee(appSerFeeReq, ctx)
}

func (s *BalanceOperationsService) WithdrawMoney(withdrawReq entity.UpdateBalanceRequest, ctx *gin.Context) (entity.UpdateBalanceResponse, error) {
	return s.repo.WithdrawMoney(withdrawReq, ctx)
}

func (s *BalanceOperationsService) Transfer(transferReq entity.TransferRequest, ctx *gin.Context) (entity.TransferResponse, error) {
	return s.repo.Transfer(transferReq, ctx)
}

func (s *BalanceOperationsService) FailedServiceFee(failedServiceFeeReq entity.StatusServiceFeeRequest, ctx *gin.Context) (entity.StatusServiceFeeResponse, error) {
	return s.repo.FailedServiceFee(failedServiceFeeReq, ctx)
}

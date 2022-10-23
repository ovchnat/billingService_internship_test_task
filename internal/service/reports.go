package service

import (
	"billingService/internal/entity"
	"billingService/internal/repository"
	"github.com/gin-gonic/gin"
)

type ReportOperationsService struct {
	repo repository.ReportOperations
}

func NewReportOperationsService(repo repository.ReportOperations) *ReportOperationsService {
	return &ReportOperationsService{repo: repo}
}

func (s *ReportOperationsService) WriteServiceMonthlyReport(serFeeReq entity.ServiceMonthlyReportReq, ctx *gin.Context) (entity.ServiceMonthlyReportResponse, error) {
	return s.repo.WriteServiceMonthlyReport(serFeeReq, ctx)
}

func (s *ReportOperationsService) GetTransactions(getTransactionsReq entity.GetTransactionsReq, ctx *gin.Context) (entity.GetTransactionsResponse, error) {
	return s.repo.GetTransactions(getTransactionsReq, ctx)
}

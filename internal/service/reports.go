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

func (s *ReportOperationsService) ServiceMonthlyReport(serFeeReq entity.ServiceMonthlyReportReq, ctx *gin.Context) (entity.ServiceMonthlyReportResponse, error) {
	return s.repo.ServiceMonthlyReport(serFeeReq, ctx)
}

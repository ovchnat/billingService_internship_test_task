package repository

import (
	"billingService/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/sqltocsv"
)

//goland:noinspection ALL
const getServicesWithSumQuery = ` SELECT service_id, SUM(invoice) 
	FROM service_log 
	WHERE status = 'Approved' AND
 	DATE(updated_at) BETWEEN $1 AND $2
	GROUP BY service_id
`

type ReportPostgres struct {
	db *sqlx.DB
}

func NewReportPostgres(db *sqlx.DB) *ReportPostgres {
	return &ReportPostgres{db: db}
}

func (r *ReportPostgres) GetTransactions(serviceReq entity.GetTransactionsReq, ctx *gin.Context) (entity.GetTransactionsResponse, error) {
	var input entity.GetTransactionsResponse

	return input, nil
}

func (r *ReportPostgres) ServiceMonthlyReport(serviceMonthlyReport entity.ServiceMonthlyReportReq, ctx *gin.Context) (entity.ServiceMonthlyReportResponse, error) {
	var input entity.ServiceMonthlyReportResponse

	rows, err := r.db.Query(getServicesWithSumQuery, serviceMonthlyReport.DateFrom, serviceMonthlyReport.DateTo)

	defer rows.Close()

	if err != nil {
		panic(err)
	}

	err = sqltocsv.WriteFile("../servicesReport.csv", rows)
	if err != nil {
		panic(err)
	}

	return input, nil
}

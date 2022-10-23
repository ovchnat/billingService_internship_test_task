package repository

import (
	"billingService/internal/entity"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/sqltocsv"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const filePathFromRoot = "../reports/"

//goland:noinspection ALL
const getServicesWithSumQuery = ` SELECT service_id, SUM(invoice) 
	FROM service_log 
	WHERE status = 'Approved' AND
 	DATE(updated_at) BETWEEN $1 AND $2
	GROUP BY service_id
`

//goland:noinspection ALL
const getTransactionByUserFromViewQuery = ` SELECT * from 
	TransactionsByAccount WHERE account_id = (SELECT id FROM accounts WHERE user_id = $1) 
	AND DATE(updated_at) BETWEEN $2 AND $3
	ORDER BY $4
	LIMIT 10 OFFSET $5;
`

type ReportPostgres struct {
	db *sqlx.DB
}

func NewReportPostgres(db *sqlx.DB) *ReportPostgres {
	return &ReportPostgres{db: db}
}

func (r *ReportPostgres) GetTransactions(getTransactionsReq entity.GetTransactionsReq, c *gin.Context) (entity.GetTransactionsResponse, error) {
	var getTransactionsResponse entity.GetTransactionsResponse
	var sortByParameter, sortByOrder string
	var idHolder int

	fail := func(err error) (entity.GetTransactionsResponse, error) {
		return getTransactionsResponse, fmt.Errorf("DepositMoney: %v", err)
	}

	switch getTransactionsReq.SortBy {
	case "date":
		sortByParameter = "created_at"
	case "amount":
		sortByParameter = "transaction_sum"
	default:

	}
	switch getTransactionsReq.SortOrder {
	case "descending":
		sortByOrder = "DESC"
	case "ascending":
		sortByOrder = "ASC"
	default:

	}
	if err := r.db.QueryRow(accByIdQuery, getTransactionsReq.UserId).Scan(&idHolder); err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("no account with that receiver id: add a new one by depositing money")
			return fail(err)
		}
		return fail(err)
	}

	rows, err := r.db.QueryContext(c, getTransactionByUserFromViewQuery,
		getTransactionsReq.UserId, getTransactionsReq.DateFrom, getTransactionsReq.DateTo,
		sortByParameter+" "+sortByOrder, (getTransactionsReq.Page-1)*10)
	if err != nil {
		return getTransactionsResponse, err
	}
	defer func(rows *sql.Rows) error {
		err := rows.Close()
		if err != nil {
			return err
		}
		return nil
	}(rows)

	fileName, newErr := TouchReportFile("userTransactionsReport")
	if newErr != nil {
		return getTransactionsResponse, err
	}
	err = sqltocsv.WriteFile(filePathFromRoot+fileName, rows)
	if err != nil {
		panic(err)
	}

	return entity.GetTransactionsResponse{
		FileLink: filePathFromRoot + fileName,
	}, nil

	return getTransactionsResponse, nil
}

func (r *ReportPostgres) WriteServiceMonthlyReport(serviceMonthlyReport entity.ServiceMonthlyReportReq, c *gin.Context) (entity.ServiceMonthlyReportResponse, error) {
	var serviceMonthlyReportResponse entity.ServiceMonthlyReportResponse

	rows, err := r.db.QueryContext(c, getServicesWithSumQuery, serviceMonthlyReport.DateFrom, serviceMonthlyReport.DateTo)

	defer func(rows *sql.Rows) error {
		err := rows.Close()
		if err != nil {
			return err
		}
		return nil
	}(rows)

	fileName, newErr := TouchReportFile("servicesReport")
	if newErr != nil {
		return serviceMonthlyReportResponse, err
	}
	err = sqltocsv.WriteFile(filePathFromRoot+fileName, rows)
	if err != nil {
		panic(err)
	}

	return entity.ServiceMonthlyReportResponse{
		FileLink: filePathFromRoot + fileName,
	}, nil
}

/* ---------------- Helper functions -------------- */

// TouchReportFile sets the file name to current date and creates a file if it doesn't exists */

func TouchReportFile(prefixString string) (string, error) {
	const dateLayout = "01-02-2006"

	fileName := prefixString + time.Now().Format(dateLayout) + ".csv"

	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	err = file.Close()
	if err != nil {
		return "", err
	}

	return fileName, nil
}

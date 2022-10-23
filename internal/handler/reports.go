package handler

import (
	"billingService/internal/entity"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

// @Summary servicesMonthly
// @Tags reports
// @Description "Get monthly report"
// @Accept json
// @Produce json
// @Param input body entity.ServiceMonthlyReportReq true "JSON object with service ID, date from and date to"
// @Success 200 {object} entity.ServiceMonthlyReportResponse
// @Failure 500 {object} errorAcc
// @Router /reports/servicesMonthly [post]
func (h *Handler) servicesMonthly(c *gin.Context) {
	var input entity.ServiceMonthlyReportReq

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.ReportOperations.WriteServiceMonthlyReport(input, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"csv-file-link": response.FileLink,
	})
}

// @Summary getCSVFile
// @Tags reports
// @Description "download CSV report file"
// @Failure 500 {object} errorAcc
// @Router /reports/{path} [get]
func (h *Handler) getCSVFile(c *gin.Context) {
	filePath := c.Param("path")

	reportFile, err := ioutil.ReadFile("../reports/" + filePath)

	if err != nil {
		panic(err)
	}
	c.Header("Content-Description", "Monthly Services Report")
	c.Header("Content-Disposition", "attachment; filename="+filePath)
	c.Data(http.StatusOK, "text/csv", reportFile)
}

// @Summary transactions
// @Tags reports
// @Description "Print user transactions log"
// @Accept json
// @Produce json
// @Param input body entity.GetTransactionsReq true "JSON object with user ID, sorting method , date from and date to"
// @Success 200 {object} entity.GetTransactionsResponse
// @Failure 500 {object} errorAcc
// @Router /reports/transactions [post]
func (h *Handler) transactions(c *gin.Context) {
	var input entity.GetTransactionsReq

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.services.ReportOperations.GetTransactions(input, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"csv-file-link": response.FileLink,
	})
}

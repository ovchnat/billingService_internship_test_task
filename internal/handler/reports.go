package handler

import (
	"billingService/internal/entity"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

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

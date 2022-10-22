package handler

import (
	"billingService/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) servicesMonthly(c *gin.Context) {
	var input entity.ServiceMonthlyReportReq

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	accResponse, err := h.services.ReportOperations.ServiceMonthlyReport(input, c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"csv-file": accResponse.ServiceId,
	})
}

func (h *Handler) transactions(c *gin.Context) {

}

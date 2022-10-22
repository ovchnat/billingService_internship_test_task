package entity

type ServiceMonthlyReportReq struct {
	ServiceId int64  `json:"service-id"`
	DateFrom  string `json:"date-from"`
	DateTo    string `json:"date-to"`
}

type ServiceMonthlyReportResponse struct {
	ServiceId int64 `json:"service-id"`
	Sum       int64 `json:"date-from"`
}

type GetTransactionsReq struct {
}

type GetTransactionsResponse struct {
}

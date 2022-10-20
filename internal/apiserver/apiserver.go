package apiserver

import (
	"billingService/internal/usecase/repository"
	"billingService/internal/usecase/service"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	services *service.BillingService
}

func NewHandler(services *service.BillingService) *Handler {
	return &Handler{services: services}
}

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start(Port string) error {

	repos := repository.NewRepo()
	services := service.NewService(repos)
	billingHandler := NewHandler(services)
	billingRouter := billingHandler.configureRoutes()

	s.httpServer = &http.Server{
		Addr:    ":" + Port,
		Handler: billingRouter,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (h *Handler) configureRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/sayHello", h.sayHello)

	accountChanges := router.Group("/accountChanges/")
	{
		accountChanges.POST("/depositMoney", h.depositMoney)
		accountChanges.POST("/withdrawMoney", h.withdrawMoney)
		accountChanges.GET("/getBalance", h.getBalance)
		accountChanges.POST("/reserveAmount", h.reserveAmount)
		accountChanges.POST("/orderConfirm", h.orderConfirm)
		accountChanges.POST("/transfer", h.transfer)

	}

	reports := router.Group("/reports")
	{
		reports.GET("/servicesMonthly", h.servicesMonthly) // csv report for the accounting
		reports.GET("/transactions", h.transactions)
	}

	return router
}

func (h *Handler) sayHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World!",
	})
}

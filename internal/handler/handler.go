package handler

import (
	_ "billingService/docs"
	"billingService/internal/repository"
	"billingService/internal/service"
	"context"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
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

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("Can't establish connection to database: %s", err.Error())
	} else {
		logrus.Println("Database connection successfully established.")
	}

	repos := repository.NewRepo(db)
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

	accountChanges := router.Group("/account/")
	{
		accountChanges.GET("/getBalance/:id", h.getBalance)
		accountChanges.POST("/depositMoney", h.depositMoney)
		accountChanges.POST("/withdrawMoney", h.withdrawMoney)
		accountChanges.POST("/reserveServiceFee", h.reserveServiceFee)
		accountChanges.POST("/approveServiceFee", h.approveOrderFee)
		accountChanges.POST("/transfer", h.transfer)
		accountChanges.POST("/failedServiceFee", h.failedServiceFee)

	}

	reports := router.Group("/reports")
	{
		reports.POST("/servicesMonthly", h.servicesMonthly) // create csv report for the accounting
		reports.GET("/:path", h.getCSVFile)
		reports.POST("/transactions", h.transactions) // create csv report for user transactions
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

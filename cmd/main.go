package main

import (
	"billingService/internal/handler"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Billing App API
// @version 1.0
// @description API server for Avito billing app

// @host localhost:8080
// @BasePath /

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("Can't initialize config: %s", err.Error())
	}

	if err := godotenv.Load("../.env"); err != nil {
		logrus.Fatalf("Can't load env variables: %s", err.Error())
	}

	srv := handler.NewServer()
	if err := srv.Start(viper.GetString("port")); err != nil {
		logrus.Fatalf("Can't get the server running: %s", err.Error())
	}
}

func initConfig() error {
	viper.SetConfigFile("../config/config.yml")
	return viper.ReadInConfig()
}

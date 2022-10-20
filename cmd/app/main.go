package main

import (
	"billingService/internal/apiserver"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("Can't initialize configs: %s", err.Error())
	}
	srv := apiserver.NewServer()
	if err := srv.Start("8080"); err != nil {
		log.Fatalf("Can't get the server running: %s", err.Error())
	}
}

func initConfig() error {
	viper.SetConfigFile("configs/config.yaml")
	return viper.ReadInConfig()
}

package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/vaidik-bajpai/medibridge/internal/handlers"
	database "github.com/vaidik-bajpai/medibridge/internal/prisma"
	"github.com/vaidik-bajpai/medibridge/internal/store"
	"go.uber.org/zap"
)

type Config struct {
	serverPort string
}

// @title           MediBridge API
// @version         1.0
// @description     Backend API for MediBridge, a medical record management system.

// @contact.name   Vaidik Bajpai
// @contact.email  codervaidik@gmail.com

// @host      localhost:8080
// @BasePath  /

// @schemes http
func main() {
	var config Config
	flag.StringVar(&config.serverPort, "sAddr", "8080", "http server address")
	flag.Parse()

	validate := validator.New()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	prismaClient, err := database.NewPrismaClient()
	if err != nil {

		panic(err)
	}
	defer prismaClient.Disconnect()

	store := store.NewStore(prismaClient)

	hdl := handlers.NewHandler(validate, logger, store)

	logger.Info("Starting the server.", zap.String("port", config.serverPort))

	err = http.ListenAndServe(fmt.Sprintf(":%s", config.serverPort), hdl.Router())
	if err != nil {
		logger.Error("error starting the server.", zap.Error(err))
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/vaidik-bajpai/medibridge/internal/handlers"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
	"github.com/vaidik-bajpai/medibridge/internal/store"
	"go.uber.org/zap"
)

type Config struct {
	serverPort string
}

func main() {
	var config Config
	flag.StringVar(&config.serverPort, "sAddr", "8080", "http server address")
	flag.Parse()

	validate := validator.New()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	prismaClient, err := NewPrismaClient()
	if err != nil {
		panic(err)
	}
	defer prismaClient.Disconnect()
	store := store.NewStore(prismaClient)

	hdl := handlers.NewHandler(validate, logger, store)

	logger.Info("Starting the server:", zap.String("port", config.serverPort))
	err = http.ListenAndServe(fmt.Sprintf(":%s", config.serverPort), hdl.Router())
	if err != nil {
		logger.Error("server crash", zap.Error(err))
	}
}

func NewPrismaClient() (*db.PrismaClient, error) {
	client := db.NewClient()
	if err := client.Connect(); err != nil {
		log.Println(err)
		return nil, err
	}

	return client, nil
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

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

	r.Post("/signup", hdl.HandleUserSignup)

	log.Println("Starting http server")
	http.ListenAndServe(fmt.Sprintf(":%s", config.serverPort), r)
}

func NewPrismaClient() (*db.PrismaClient, error) {
	client := db.NewClient()
	if err := client.Connect(); err != nil {
		return nil, err
	}

	return client, nil
}

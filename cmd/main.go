package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/v1", func(r chi.Router) {
		r.Post("/signup", hdl.HandleUserSignup)
		r.Post("/signin", hdl.HandleUserLogin)
		r.Route("/patient", func(r chi.Router) {
			r.Use(hdl.RequireAuth)
			r.Post("/", hdl.HandleRegisterPatient)
			r.Route("/{patientID}", func(r chi.Router) {
				r.Put("/", hdl.HandleUpdatePatientDetails)
				r.Delete("/", hdl.HandleDeletePatientDetails)

				r.Post("/diagnoses", hdl.HandleAddDiagnoses)
			})

			r.Route("/{diagnosesID}", func(r chi.Router) {
				r.Put("/", hdl.HandleUpdateDiagnoses)
				r.Delete("/", hdl.HandleDeleteDiagnoses)
			})
		})
	})

	logger.Info("Starting the server:", zap.String("port", config.serverPort))
	http.ListenAndServe(fmt.Sprintf(":%s", config.serverPort), r)
}

func NewPrismaClient() (*db.PrismaClient, error) {
	client := db.NewClient()
	if err := client.Connect(); err != nil {
		return nil, err
	}

	return client, nil
}

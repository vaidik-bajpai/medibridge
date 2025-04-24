// @title MediBridge API
// @version 1.0
// @description This is the API documentation for the MediBridge backend.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.email support@medibridge.com

// @host localhost:8080
// @BasePath /
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"                 // For input validation
	_ "github.com/joho/godotenv/autoload"                    // Automatically loads environment variables from .env
	"github.com/vaidik-bajpai/medibridge/internal/handlers"  // Importing custom handlers
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db" // Importing Prisma DB client
	"github.com/vaidik-bajpai/medibridge/internal/store"     // Importing store for business logic
	"go.uber.org/zap"                                        // For structured logging
)

// Config struct holds the server configuration, such as the server port
type Config struct {
	serverPort string
}

func main() {
	// Initialize configuration using command-line flags
	var config Config
	flag.StringVar(&config.serverPort, "sAddr", "8080", "http server address") // Default to port 8080
	flag.Parse()

	// Create a new validator instance for validating inputs
	validate := validator.New()

	// Initialize a logger with production-level configuration
	logger, _ := zap.NewProduction()
	defer logger.Sync() // Ensure logs are flushed when the program ends

	// Initialize Prisma client for database interactions
	prismaClient, err := NewPrismaClient()
	if err != nil {
		// Log any errors and panic if client initialization fails
		panic(err)
	}
	defer prismaClient.Disconnect() // Ensure the client disconnects after usage

	// Create a new store instance that interacts with the Prisma client
	store := store.NewStore(prismaClient)

	// Initialize handler with validation, logger, and store
	hdl := handlers.NewHandler(validate, logger, store)

	// Log server start message with the chosen port
	logger.Info("Starting the server:", zap.String("port", config.serverPort))

	// Start the HTTP server and listen for incoming requests
	// @Summary Starts the HTTP server
	// @Description Starts the HTTP server to handle requests.
	// @Tags server
	// @Accept json
	// @Produce json
	// @Param port query string false "Port to listen on"
	// @Success 200 {string} string "Server started"
	// @Failure 500 {string} string "Server crash"
	// @Router /start [get]
	err = http.ListenAndServe(fmt.Sprintf(":%s", config.serverPort), hdl.Router())
	if err != nil {
		// If the server crashes, log the error
		logger.Error("server crash", zap.Error(err))
	}
}

// NewPrismaClient initializes and returns a new Prisma client, connected to the database
// @Summary Initializes the Prisma client
// @Description Initializes the Prisma client and connects to the database.
// @Tags database
// @Success 200 {string} string "Prisma client initialized"
// @Failure 500 {string} string "Failed to connect to the database"
func NewPrismaClient() (*db.PrismaClient, error) {
	// Create a new Prisma client instance
	client := db.NewClient()

	// Attempt to connect to the Prisma client
	if err := client.Connect(); err != nil {
		// Log any connection errors and return the error
		log.Println(err)
		return nil, err
	}

	// Return the successfully connected Prisma client
	return client, nil
}

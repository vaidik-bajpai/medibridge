package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
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
	r.Get("/signup", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Sign up form"))
	})

	http.ListenAndServe(fmt.Sprintf(":%s", config.serverPort), r)
}

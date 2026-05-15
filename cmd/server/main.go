package main

import (
	"awesomeProject1/internal"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	addr := os.Getenv("SERVER_ADDR")
	conf := os.Getenv("CONF")

	repository := internal.NewConfRepository(conf)
	service := internal.NewService(repository)
	API := internal.NewAPI(service)

	srv := &http.Server{
		Addr:              addr,
		Handler:           API.Routes(),
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Println("HTTP server on", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

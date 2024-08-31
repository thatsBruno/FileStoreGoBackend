package main

import (
	"go-api/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// database actions
	db.InitDatabase()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /user/{id}", UserHandlers)
	mux.HandleFunc("POST /user", UserHandlers)
	mux.HandleFunc("DELETE /user/{id}", UserHandlers)

	mux.HandleFunc("GET /files/{ownerid}", FileHandlers)
	mux.HandleFunc("POST /file", FileHandlers)
	mux.HandleFunc("DELETE /file/{fileid}", FileHandlers)

	server := &http.Server{
		Addr:    ":4000",
		Handler: mux,
	}

	// channel to listen for interrupt signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	go func() {
		log.Println("Starting server on :4000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :4000: %v\n", err)
		}
	}()

	<-stopChan

	log.Println("Shutting down server...")
	server.SetKeepAlivesEnabled(false)
	if err := server.Close(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	time.Sleep(5 * time.Second)

	log.Println("Server exiting")
}

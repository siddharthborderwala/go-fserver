package main

import (
	"context"
	"log"
	"microserver/files"
	"microserver/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const BIND_ADDRESS = ":9090"
const BASE_PATH string = "./imagestore"
const SIZE int = 1024 * 1024 * 5

func main() {
	logger := log.Default()

	// instantiate storage
	store, err := files.NewLocal(BASE_PATH, SIZE)
	if err != nil {
		logger.Fatal("unable to create storage", "error", err)
	}

	// create the handlers
	fh := handlers.NewFiles(store, logger)
	gzip := handlers.NewGzipMiddleware()

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// post method handler
	ph := sm.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+.[a-z]{3}}", fh.UploadRest)
	ph.HandleFunc("/", fh.UploadMultipart)

	// get method handler
	gh := sm.Methods(http.MethodGet).Subrouter()
	gh.Use(gzip.Middleware)
	gh.Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z]+.[a-z]{3}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir(BASE_PATH))),
	)

	s := http.Server{
		Addr:         BIND_ADDRESS,
		Handler:      ch(sm),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Println("Starting server", "address", BIND_ADDRESS)
		err := s.ListenAndServe()
		if err != nil {
			logger.Fatal("unable to start server", "error", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// block until a signal is received
	sig := <-sigChan
	logger.Println("Shutting down server, got", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-service/api/internal/handler"
	"auth-service/api/internal/server"
)

func main() {
	ds, err := server.InitDS()

	if err != nil {
		log.Fatalf("Unable to initialize data sources %v", err)
	}

	router, err := server.Inject(ds)
	if err != nil {
		log.Fatalf("Unable to inject data sources %v", err)
	}

	handler.NewHandler(&handler.Config{
		R: router,
	})

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialized server: %v\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown data sources
	if err := ds.Close(); err != nil {
		log.Fatalf("A problem occurred gracefully shutting down data sources: %v\n", err)
	}

	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %v\n", err)
	}
}

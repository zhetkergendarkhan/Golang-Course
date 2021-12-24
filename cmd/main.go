package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"shop/internal/transport/http"
	"syscall"
	"time"
)

func main() {

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	var errCh chan error
	go gracefulShutdown(errCh)

	go http.StartServer(errCh)
	log.Print("Service started on port: 8000")
	err := <-errCh

	log.Println("Gracefully stop service")

	<-time.NewTicker(5 * time.Second).C

	log.Printf("Service terminated: %v\n", err)
}

func gracefulShutdown(errCh chan<- error) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	errCh <- fmt.Errorf("service is shutting down. Signal %d", <-sigCh)
}

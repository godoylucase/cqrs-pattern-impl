package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godoylucase/cqrs-pattern-impl/business"
	"github.com/godoylucase/cqrs-pattern-impl/external/api"
	"github.com/godoylucase/cqrs-pattern-impl/internal/db"
	"github.com/godoylucase/cqrs-pattern-impl/internal/repository"
	"github.com/sirupsen/logrus"
)

const port = 8081

func main() {
	done := make(chan interface{})
	defer close(done)

	mongoConn := db.NewMongoDBConn()
	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := mongoConn.Connect(dbCtx)
	if err != nil {
		logrus.Fatal(err)
	}
	defer mongoConn.Disconnect(dbCtx)

	eb := business.NewEventBroker(done)

	ar, err := repository.NewArticleRepository(mongoConn)
	if err != nil {
		logrus.Fatal(err)
	}

	as := business.NewArticleService(ar, eb)

	h := &api.Handler{
		As: as,
	}
	api := api.Configure(h, gin.Default())

	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		Handler:      api,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Printf("Starting HTTP Server. Listening at %q \n", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("%v", err)
		} else {
			log.Println("Server closed!")
		}
	}()

	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Interrupt, syscall.SIGTERM)
	sig := <-sigquit
	done <- 1
	log.Printf("caught sig: %+v\n", sig)
	log.Printf("Gracefully shutting down server...\n")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("Unable to shut down server: %v\n", err)
	} else {
		log.Println("Server stopped")
	}
}

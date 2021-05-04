package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godoylucase/read_tags/business"
	"github.com/godoylucase/read_tags/business/news"
	"github.com/godoylucase/read_tags/external/api"
	"github.com/godoylucase/read_tags/external/event"
	"github.com/godoylucase/read_tags/internal/db"
	"github.com/godoylucase/read_tags/internal/repository"
)

const port = 8082

func main() {
	cassandraClient, err := db.Cassandra()
	if err != nil {
		panic(err)
	}
	defer cassandraClient.Close()

	artRepo := repository.NewArticle(cassandraClient)
	resolverFactory := news.NewResolverFactory(artRepo)

	articleResolver, err := resolverFactory.Get(event.ARTICLE)
	if err != nil {
		panic(err)
	}

	done := make(chan interface{})

	articleConsumer, err := event.NewConsumer(done)
	if err != nil {
		panic(err)
	}

	articleProcessor := news.NewProcessor(done, event.ARTICLE, articleResolver, articleConsumer)

	// run event processing in background
	go func() {
		if err := articleProcessor.Run(); err != nil {
			panic(err)
		}
	}()

	repo := repository.NewArticle(cassandraClient)

	qs := business.NewQueryService(repo)

	h := api.NewHandler(qs)
	apiHandler := api.Configure(h, gin.Default())

	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		Handler:      apiHandler,
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
	close(done)
	log.Printf("caught sig: %+v\n", sig)
	log.Printf("Gracefully shutting down server...\n")

}

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/godoylucase/read_tags/business/news"
	"github.com/godoylucase/read_tags/external/event"
)

func main() {

	artRepo := &news.ARMock{}
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

	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Interrupt, syscall.SIGTERM)
	sig := <-sigquit
	close(done)
	log.Printf("caught sig: %+v\n", sig)
	log.Printf("Gracefully shutting down server...\n")

}

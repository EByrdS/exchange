package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	engineserver "exchange/engine/server"
)

func killSignal() chan os.Signal {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	return sigc
}

func main() {
	fmt.Println("Welcome to the engine.")

	markets := []engineserver.MarketSymbol{
		{Base: "DOLS", Trade: "MEEM"},
	}
	engine, err := engineserver.NewEngine(markets)
	if err != nil {
		log.Fatalf("Error creating engine: %v\n", err)
	}

	defer engine.CloseKafka()

	ctx := context.Background()
	fmt.Println("Engine starting to stream messages...")
	engine.Stream(ctx)

	errChan := make(chan error, 1)
	go func() {
		fmt.Println("Engine starting to listen for Kafka messages...")
		if err := engine.Listen(ctx); err != nil {
			errChan <- err
		}
	}()

	sigc := killSignal()
	<-sigc
	fmt.Println("Shutting down gracefully...")
}

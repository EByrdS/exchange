package main

import (
	"context"
	engineserver "exchange/engine/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/twmb/franz-go/pkg/kgo"
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

type Printer struct {
	kafka *kgo.Client
}

func NewPrinter(markets []engineserver.MarketSymbol) (*Printer, error) {
	topics := []string{}
	for _, market := range markets {
		topics = append(topics,
			market.Topic()+".orders",
			market.Topic()+".volumes",
			market.Topic()+".matches",
		)
	}

	cl, err := kgo.NewClient(
		kgo.SeedBrokers("localhost:9092"),
		kgo.ConsumeTopics(topics...),
		kgo.ConsumerGroup("printer"),
	)
	if err != nil {
		return nil, err
	}

	p := &Printer{kafka: cl}

	return p, nil
}

func main() {
	fmt.Println("Welcome to the recorder (printer)")

	markets := []engineserver.MarketSymbol{
		{Base: "DOLS", Trade: "MEEM"},
	}
	printer, err := NewPrinter(markets)
	if err != nil {
		log.Fatalf("Error creating printer: %v", err)
	}

	ctx := context.Background()
	printer.Listen(ctx)
	fmt.Println("Message listening ready")

	sigc := killSignal()

	<-sigc
	fmt.Println("Shutting down gracefully...")
}

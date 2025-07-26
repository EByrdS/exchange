package engineserver

import (
	"context"
	"exchange/engine/market"
	"log"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Engine is a struct that coordinates events between the order books of different
// trading pairs, and the Kafka client.
type Engine struct {
	// A sync Map with all trading pairs available, from symbol topics to market
	pairs sync.Map

	// The market symbols available
	marketSymbols []MarketSymbol

	// A map from symbol topics to order event channels
	orderEventsChans map[string]chan *market.OrderEvent

	// A map from symbol topics to volume event channels
	volumeEventsChans map[string]chan *market.VolumeEvent

	// A map from symbol topics to match event channels
	matchEventsChans map[string]chan *market.MatchEvent

	// The Kafka client
	kafka *kgo.Client
}

// NewEngine creates an engine with initialized channels and a kafka client.
func NewEngine(markets []MarketSymbol) (*Engine, error) {
	topics := []string{}
	for _, market := range markets {
		topics = append(topics, market.Topic())
	}

	cl, err := kgo.NewClient(
		kgo.SeedBrokers("localhost:9092"),
		kgo.ConsumeTopics(topics...),
		kgo.ConsumerGroup("engine"),
	)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	ctxTime, cancel := context.WithTimeout(ctx, 2*time.Second)
	err = cl.Ping(ctxTime)
	cancel()
	if err != nil {
		log.Fatal("Ping unsuccessful, kafka client down?")
	} else {
		log.Print("Ping successful, kafka client is up!")
	}

	e := &Engine{
		marketSymbols:     markets,
		orderEventsChans:  map[string]chan *market.OrderEvent{},
		volumeEventsChans: map[string]chan *market.VolumeEvent{},
		matchEventsChans:  map[string]chan *market.MatchEvent{},
		kafka:             cl,
	}

	for _, market := range markets {
		e.addMarket(market)
	}

	return e, nil
}

func (e *Engine) addMarket(ms MarketSymbol) {
	orderEventsChan := make(chan *market.OrderEvent, 20)
	volumeEventsChan := make(chan *market.VolumeEvent, 20)
	matchEventsChan := make(chan *market.MatchEvent, 20)

	e.orderEventsChans[ms.Topic()] = orderEventsChan
	e.volumeEventsChans[ms.Topic()] = volumeEventsChan
	e.matchEventsChans[ms.Topic()] = matchEventsChan

	m := market.New(ms.Name(), orderEventsChan, volumeEventsChan, matchEventsChan)

	e.pairs.Store(ms.Topic(), m)
}

func (e *Engine) CloseKafka() {
	if e.kafka != nil {
		e.kafka.Close()
	}
}

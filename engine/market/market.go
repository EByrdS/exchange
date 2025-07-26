package market

import (
	"exchange/engine/order"
	"exchange/engine/orderbook"
	"exchange/engine/orderbook/rbtree"
	"sync"
	"time"
)

// Market has both order books of a trading market, and is responsible
// for firing the corresponding events.
type Market struct {
	// This market pair name.
	pair string

	// The buy side order book.
	buyBook *orderbook.OrderBook

	// The sell side order book.
	sellBook *orderbook.OrderBook

	// Events that communicate individual order's lifetime in the market.
	orderEvents chan<- *OrderEvent

	// Events triggered when two orders are matched.
	matchEvents chan<- *MatchEvent
}

func New(pair string, orderEvents chan<- *OrderEvent, volumeEvents chan<- *VolumeEvent, matchEvents chan<- *MatchEvent) *Market {
	pool := &sync.Pool{
		New: func() any {
			return rbtree.NewNode()
		},
	}

	m := &Market{
		pair:        pair,
		orderEvents: orderEvents,
		matchEvents: matchEvents,
		buyBook: orderbook.New(order.OrderBuy, pool, func(price uint64, volume uint64) {
			volumeEvents <- &VolumeEvent{Pair: pair, Side: order.OrderBuy, Price: price, Volume: volume, Timestamp: time.Now()}
		}),
		sellBook: orderbook.New(order.OrderSell, pool, func(price uint64, volume uint64) {
			volumeEvents <- &VolumeEvent{Pair: pair, Side: order.OrderSell, Price: price, Volume: volume, Timestamp: time.Now()}
		}),
	}

	return m
}

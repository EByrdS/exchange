package market_test

import (
	"exchange/engine/market"
	"exchange/engine/order"
	"fmt"
	"testing"
)

// eventsTrackers is a test helper struct that creates buffered channels
// and reads data until there are no messages left. It is meant to be used in
// the same goroutine as the tests, and
type eventsTracker struct {
	volumeEventsChan chan *market.VolumeEvent
	orderEventsChan  chan *market.OrderEvent
	matchEventsChan  chan *market.MatchEvent

	volumeEvents []*market.VolumeEvent
	orderEvents  []*market.OrderEvent
	matchEvents  []*market.MatchEvent
}

// initializes an events tracker with the corresponding buffer size. The tests
// should have a good estimate on how many messages will be received, otherwise
// they will block execution if the channels run out of buffered space.
func newEventsTracker(bufferSize int) *eventsTracker {
	return &eventsTracker{
		volumeEventsChan: make(chan *market.VolumeEvent, bufferSize),
		orderEventsChan:  make(chan *market.OrderEvent, bufferSize),
		matchEventsChan:  make(chan *market.MatchEvent, bufferSize),
	}
}

// flush will read from every buffered channel until there are no messages left
// and stores the messages in their corresponding slice. The channels are still
// receiving after calling this function.
func (e *eventsTracker) flush() {
	for {
		if len(e.volumeEventsChan) == 0 {
			break
		}

		ev, ok := <-e.volumeEventsChan
		if !ok {
			break
		}

		e.volumeEvents = append(e.volumeEvents, ev)
	}

	for {
		if len(e.orderEventsChan) == 0 {
			break
		}

		ev, ok := <-e.orderEventsChan
		if !ok {
			break
		}

		e.orderEvents = append(e.orderEvents, ev)
	}

	for {
		if len(e.matchEventsChan) == 0 {
			break
		}

		ev, ok := <-e.matchEventsChan
		if !ok {
			break
		}

		e.matchEvents = append(e.matchEvents, ev)
	}
}

// reset clears the channel queues and then clearss the stored events
// resulting in a clean state ready for testing
func (e *eventsTracker) reset() {
	e.flush()

	e.volumeEvents = []*market.VolumeEvent{}
	e.orderEvents = []*market.OrderEvent{}
	e.matchEvents = []*market.MatchEvent{}
}

func (e *eventsTracker) ignoreAll() {
	go func() {
		for {
			<-e.volumeEventsChan
		}
	}()

	go func() {
		for {
			<-e.orderEventsChan
		}
	}()

	go func() {
		for {
			<-e.matchEventsChan
		}
	}()
}

func Benchmark_PriceDeletionAndInsertion(b *testing.B) {
	// On this test we expect allocations to happen only for the triggered events
	// not for the creation or price nodes.

	// After examining this test, we conclude that indeed, recycling tree nodes
	// allocates little memory, but creating price levels is still allocating
	// as usual. This is a design flaw, a tree node contains a price level, and
	// the tree node is recycled but the price level is not.
	//
	// The price level contains complex structures: a doubly linked list and a
	// map, both of which allocate memory every time. One solution could be to
	// create a price level sync Pool, but each struct would have to initialize
	// its own complex fields anyways. We could try to go deeper and create either
	// a pool that keeps initialized price levels intact, and recycles them, or
	// creating a different pool for each of its fields in turn. Still, any
	// solution would need major refactoring, so I am choosing not to do that
	// in favor of advancing the project.
	//
	// Getting to this place, observing the result and learning this has been
	// already invaluable.

	tracker := newEventsTracker(10)
	tracker.ignoreAll()

	testCases := []struct {
		name  string
		depth uint64
	}{
		{
			name:  "ten_prices_each_side",
			depth: 10,
		},
		{
			name:  "hundred_prices_each_side",
			depth: 100,
		},
		{
			name:  "thousand_prices_each_side",
			depth: 1_000,
		},
		{
			name:  "ten_thousand_prices_each_side",
			depth: 10_000,
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			pair := "USD/BTC"
			sellBoundary := uint64(1_000_010)
			buyBoundary := uint64(1_000_005)
			m := market.New(pair, tracker.orderEventsChan, tracker.volumeEventsChan, tracker.matchEventsChan)

			buyOrders, sellOrders := []*order.Order{}, []*order.Order{}
			for i := range tc.depth {
				buyOrder := &order.Order{ID: fmt.Sprintf("buy-%d", i), Pair: pair, Price: buyBoundary - i, Side: order.OrderBuy, Volume: 1}
				sellOrder := &order.Order{ID: fmt.Sprintf("sell-%d", i), Pair: pair, Price: sellBoundary + i, Side: order.OrderSell, Volume: 1}

				buyOrders = append(buyOrders, buyOrder)
				sellOrders = append(sellOrders, sellOrder)

				// Start only with all the buy orders
				m.InsertMakerOrder(buyOrder)
			}

			b.ResetTimer()

			for range b.N {
				for _, o := range buyOrders {
					// Delete buy orders to fill the price node pool
					m.Cancel(o)
				}

				for _, o := range sellOrders {
					// Insert the other side so that they use the price node pool
					m.InsertMakerOrder(o)
				}

				for _, o := range sellOrders {
					// Delete order to re-fill the price node pool
					m.Cancel(o)
				}

				for _, o := range buyOrders {
					// Repopulate the buy side to get price nodes from the sync pool and
					// come back to the initial state
					m.InsertMakerOrder(o)
				}
			}
		})
	}
}

package market

import (
	"exchange/engine/order"
	"exchange/engine/orderbook"
	"fmt"
	"time"
)

// MatchTakerOrder will take as much volume as possible from the corresponding
// maker side until the volume of the taker order is fulfilled.
func (m *Market) MatchTakerOrder(o *order.Order) error {
	if err := m.validateTakerOrder(o); err != nil {
		return fmt.Errorf("match taker order: %w", err)
	}

	makerBook := m.sellBook
	if o.Side == order.OrderSell {
		makerBook = m.buyBook
	}

	m.matchTakerOrder(o, makerBook)

	return nil
}

func (m *Market) matchTakerOrder(o *order.Order, makerBook *orderbook.OrderBook) {
	txnTime := time.Now()

	matches, missingVolume := makerBook.MatchAndExtract(o.Volume)

	if missingVolume > 0 {
		m.orderEvents <- &OrderEvent{Type: TakerOrderUnfulfilled, OrderID: o.ID, Timestamp: txnTime}
	}

	for i, match := range matches {
		takerMatchType := order.OrderPartiallyFulfilled
		if i == len(matches)-1 && missingVolume == 0 {
			takerMatchType = order.OrderFulfilled
		}

		m.matchEvents <- &MatchEvent{
			Pair:            m.pair,
			TakerOrderID:    o.ID,
			TakerMatchType:  takerMatchType,
			MakerOrderID:    match.MakerOrder.ID,
			MakerMatchType:  match.Type,
			MatchedVolume:   match.VolumeTaken,
			SettlementPrice: match.MakerOrder.Price,
			Timestamp:       txnTime,
		}
	}
}

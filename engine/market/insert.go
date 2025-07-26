package market

import (
	"exchange/engine/order"
	"exchange/engine/orderbook"
	"fmt"
	"time"
)

// InsertMakerOrder places a maker order in its corresponding side. If the order
// would cross the market boundary it will treat it as a taker order.
func (m *Market) InsertMakerOrder(o *order.Order) error {
	if err := m.validateOrder(o); err != nil {
		return fmt.Errorf("InsertMakerOrder: %w", err)
	}

	var makerBook *orderbook.OrderBook
	if o.Side == order.OrderBuy {
		if headPrice := m.sellBook.HeadPrice(); headPrice != 0 && o.Price >= headPrice {
			m.matchTakerOrder(o, m.sellBook)
			return nil
		}

		makerBook = m.buyBook
	} else {
		if headPrice := m.buyBook.HeadPrice(); headPrice != 0 && o.Price <= headPrice {
			m.matchTakerOrder(o, m.buyBook)
			return nil
		}

		makerBook = m.sellBook
	}

	if err := makerBook.Insert(o); err != nil {
		m.orderEvents <- &OrderEvent{Type: OrderRejected, OrderID: o.ID, Timestamp: time.Now()}
		return err
	}

	m.orderEvents <- &OrderEvent{Type: MakerOrderInserted, OrderID: o.ID, Timestamp: time.Now()}
	return nil
}

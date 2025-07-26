package market

import (
	"exchange/engine/order"
	"time"
)

// Cancel removes an order from its corresponding book.
func (m *Market) Cancel(o *order.Order) error {
	if err := m.validateOrder(o); err != nil {
		return err
	}

	book := m.buyBook
	if o.Side == order.OrderSell {
		book = m.sellBook
	}

	if err := book.Delete(o); err != nil {
		return err
	}

	m.orderEvents <- &OrderEvent{Type: OrderCancelled, OrderID: o.ID, Timestamp: time.Now()}
	return nil
}

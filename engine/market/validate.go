package market

import (
	"exchange/engine/order"
	"fmt"
	"time"
)

func (m *Market) validateOrder(o *order.Order) error {
	if err := m.validateTakerOrder(o); err != nil {
		return err
	}

	if o.Price <= 0 {
		m.orderEvents <- &OrderEvent{Type: OrderRejected, OrderID: o.ID, Timestamp: time.Now()}
		return fmt.Errorf("market %q, order %q, negative or zero price %d: %w", m.pair, o.ID, o.Price, InvalidOrderErr)
	}

	return nil
}

func (m *Market) validateTakerOrder(o *order.Order) error {
	if o == nil {
		return fmt.Errorf("market %q, nil order: %w", m.pair, InvalidOrderErr)
	}

	if o.ID == "" {
		m.orderEvents <- &OrderEvent{Type: OrderRejected, Timestamp: time.Now()}
		return fmt.Errorf("market %q, no order ID: %v: %w", m.pair, o, InvalidOrderErr)
	}

	if o.Pair != m.pair {
		m.orderEvents <- &OrderEvent{Type: OrderRejected, OrderID: o.ID, Timestamp: time.Now()}
		return fmt.Errorf("market %q, order %q, different pair %q: %w", m.pair, o.ID, o.Pair, InvalidOrderErr)
	}

	if o.Volume <= 0 {
		m.orderEvents <- &OrderEvent{Type: OrderRejected, OrderID: o.ID, Timestamp: time.Now()}
		return fmt.Errorf("market %q, order %q, negative or zero volume %d: %w", m.pair, o.ID, o.Volume, InvalidOrderErr)
	}

	return nil
}

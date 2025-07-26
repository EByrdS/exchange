package pricelevel

import (
	"exchange/engine/order"
	"fmt"
)

// Remove receives an orderID and removes it from the queue, updating the
// volume accordingly.
//
// O(1).
func (p *PriceLevel) Remove(orderID string) error {
	elem, ok := p.orderMap[orderID]
	if !ok {
		return fmt.Errorf("PriceLevel unknown order ID %s", orderID)
	}

	order := elem.Value.(*order.Order)
	startLen := p.list.Len() // O(1)

	p.list.Remove(elem)
	if startLen == p.list.Len() {
		return fmt.Errorf("PriceLevel.Remove: corrupted state, cannot remove element with order id %s", order.ID)
	}

	delete(p.orderMap, orderID)

	p.volume -= order.Volume

	return nil
}

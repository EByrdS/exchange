package pricelevel

import (
	"fmt"

	"exchange/engine/order"
)

func (p *PriceLevel) insert(o *order.Order) {
	p.volume += o.Volume
	elem := p.list.PushBack(o)
	p.orderMap[o.ID] = elem
}

// Push back adds an order to be processed last and updates the Volume.
//
// O(1).
func (p *PriceLevel) Insert(order *order.Order) error {
	if _, ok := p.orderMap[order.ID]; ok {
		return fmt.Errorf("PriceLevel received duplicate order ID %s", order.ID)
	}

	p.insert(order)

	return nil
}

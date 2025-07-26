package pricelevel

import (
	"container/list"
	"exchange/engine/order"
)

// The PriceLevel receives orders and processes them FIFO, keeping track of the
// total Volume.
// Additions and deletions are O(1).
type PriceLevel struct {
	// The up-to-date volume of this whole price level
	volume uint64

	// The doubly-linked-list to process orders FIFO
	// Length: O(1)
	// Read front: O(1)
	// Read back: O(1)
	// Push back: O(1)
	// Remove element: O(1)
	// Finding an element: O(n) <- this is why we keep an orderMap
	list *list.List

	// A map from order ID to list elements for O(1) deletions
	orderMap map[string]*list.Element
}

// Volume returns the current available volume of this level.
//
// O(1)
func (p *PriceLevel) Volume() uint64 {
	return p.volume
}

// Front returns the order that is first in the processing order in O(1).
func (p *PriceLevel) Front() *order.Order {
	elem := p.list.Front()
	if elem == nil {
		return nil
	}

	return elem.Value.(*order.Order)
}

func New() *PriceLevel {
	return &PriceLevel{
		list:     list.New(),
		orderMap: make(map[string]*list.Element),
	}
}

func (p *PriceLevel) Reset() {
	if p == nil {
		return
	}

	p.volume = 0
	p.list.Init()
	for k := range p.orderMap {
		delete(p.orderMap, k)
	}
}

package orderbook

import (
	"fmt"

	"exchange/engine/order"
	"exchange/engine/orderbook/rbtree"
)

// O(log n)
func (b *OrderBook) deletePriceNode(node *rbtree.Node) {
	delete(b.priceMap, node.Price)
	b.priceTree.DeleteNode(node) // O(log n)
}

// Delete identifies the price node where the given order is, and removes it.
//
// O(log n):
// O(1) if the price node still has orders.
// O(log n) if the price node becomes empty and has to be removed.
func (b *OrderBook) Delete(o *order.Order) error {
	if o.Side != b.side {
		return fmt.Errorf("OrderBook.Delete(%q) different sides %v!=%v", o.ID, b.side, o.Side)
	}

	priceNode, exists := b.priceMap[o.Price] // O(1)
	if !exists {
		return fmt.Errorf("OrderBook.Delete(%q) price node %d does not exist", o.ID, o.Price)
	}

	if err := priceNode.Orders.Remove(o.ID); err != nil { // O(1)
		return fmt.Errorf("OrderBook.Delete(%q) failed to remove: %w", o.ID, err)
	}

	b.volumeUpdateCallback(priceNode.Price, priceNode.Orders.Volume())
	if priceNode.Orders.Volume() == 0 {
		b.deletePriceNode(priceNode) // O(log n)
	}

	return nil
}

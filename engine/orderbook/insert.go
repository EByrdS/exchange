package orderbook

import (
	"fmt"

	"exchange/engine/order"
	"exchange/engine/orderbook/rbtree"
)

// O(log n)
func (b *OrderBook) insertPriceNode(price uint64) *rbtree.Node {
	priceNode := b.priceTree.Insert(price) // O(log n)
	b.priceMap[price] = priceNode

	return priceNode
}

// Insert colocates an order in its correct price node, creating the node if it
// doesn't exist.
//
// O(log n):
// O(1) if the price already exists
// O(log n) to insert a new price node
func (b *OrderBook) Insert(o *order.Order) error {
	if o.Side != b.side {
		return fmt.Errorf("OrderBook.Insert(%q) different sides %v!=%v", o.ID, b.side, o.Side)
	}

	priceNode, exists := b.priceMap[o.Price] // O(1)
	if !exists {
		priceNode = b.insertPriceNode(o.Price) // O(log n)
	}

	if err := priceNode.Orders.Insert(o); err != nil { // O(1)
		return err
	}

	b.volumeUpdateCallback(o.Price, priceNode.Orders.Volume())
	return nil
}

package orderbook

import (
	"exchange/engine/order"
	"exchange/engine/orderbook/rbtree"
	"sync"
)

type OrderBook struct {
	// If this book is for buy or sell orders
	side order.OrderSide

	// The binary search tree to optimize price discovery
	// Searches in O(log n). Inserts in O(log n)
	priceTree *rbtree.Tree

	// A map of price to tree node to know if a price exists
	// Reads in O(1). Writes in O(1)
	priceMap map[uint64]*rbtree.Node

	// A function to call when there is a change in volume
	volumeUpdateCallback func(price uint64, volume uint64)
}

func New(side order.OrderSide, nodePool *sync.Pool, volumeUpdateCallback func(uint64, uint64)) *OrderBook {
	treeOrientation := rbtree.MinFirst
	if side == order.OrderBuy {
		treeOrientation = rbtree.MaxFirst
	}

	return &OrderBook{
		side:                 side,
		priceTree:            rbtree.NewTree(treeOrientation, nodePool),
		priceMap:             make(map[uint64]*rbtree.Node),
		volumeUpdateCallback: volumeUpdateCallback,
	}
}

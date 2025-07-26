package orderbook

import (
	"exchange/engine/order"
)

// MatchAndExtract will return all the orders needed to fill the required volume,
// and returns the volume that could not be extracted.
//
// O(n)
func (b *OrderBook) MatchAndExtract(volume uint64) ([]*order.Match, uint64) {
	totalMatches := make([]*order.Match, 0, 10)

	var matches []*order.Match
	for volume > 0 {
		head := b.priceTree.Head() // O(1)
		if head == nil {
			break
		}

		matches, volume = head.Orders.MatchAndExtract(volume) // O(n)
		b.volumeUpdateCallback(head.Price, head.Orders.Volume())

		if head.Orders.Volume() == 0 {
			b.deletePriceNode(head) // O(log n)
		}
		totalMatches = append(totalMatches, matches...)
	}

	return totalMatches, volume
}

package orderbook

// Snapshot returns an up-to-date map of [price] -> volume
//
// O(n)
func (o *OrderBook) Snapshot() map[uint64]uint64 {
	volumes := make(map[uint64]uint64, len(o.priceMap))
	for price, priceNode := range o.priceMap {
		volumes[price] = priceNode.Orders.Volume()
	}

	return volumes
}

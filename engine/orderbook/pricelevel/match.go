package pricelevel

import (
	"exchange/engine/order"
)

// Extract will return all the orders needed to fill the required volume,
// and any unmatched volume that could not be extracted from this level.
//
// O(n)
func (b *PriceLevel) MatchAndExtract(volume uint64) ([]*order.Match, uint64) {
	matches := make([]*order.Match, 0, 10)

	// O(n)
	for {
		if volume == 0 {
			break
		}

		elem := b.list.Front()
		if elem == nil {
			break
		}

		o := elem.Value.(*order.Order)
		if volume >= o.Volume {
			volume -= o.Volume
			b.volume -= o.Volume

			b.list.Remove(elem) // removing the pointer to Next
			delete(b.orderMap, o.ID)

			matches = append(matches, &order.Match{
				Type:        order.OrderFulfilled,
				MakerOrder:  o,
				VolumeTaken: o.Volume,
			})

			o.Volume = 0
		} else {
			o.Volume -= volume
			b.volume -= volume

			matches = append(matches, &order.Match{
				Type:        order.OrderPartiallyFulfilled,
				MakerOrder:  o,
				VolumeTaken: volume,
			})
			return matches, 0
		}
	}

	return matches, volume
}

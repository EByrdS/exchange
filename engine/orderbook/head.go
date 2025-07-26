package orderbook

// HeadPrice is the price at the head of the orderbook. For a sell-side book
// this will be the lowest price, and for a buy-side book this will be the
// highest price. Returns 0 if there is no head price.
//
// O(1)
func (o OrderBook) HeadPrice() uint64 {
	head := o.priceTree.Head()
	if head == nil {
		return 0
	}

	return head.Price
}

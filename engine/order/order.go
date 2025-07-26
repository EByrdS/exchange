package order

import "fmt"

type OrderSide int

const (
	OrderBuy OrderSide = iota
	OrderSell
)

// Order contains the information needed for the orderbook to operate.
// It does not contain meta-information, to keep the orderbook lightweight.
type Order struct {
	// A unique identifier for this order.
	ID string

	// The market pair name the order is transacting in.
	Pair string

	// Whether the order is on the buy or sell side.
	Side OrderSide

	// The price specified by the user if this is a limit order (maker).
	// A market order (taker) will have a price of 0.
	Price uint64

	// The asset quantity to be exchanged, NOT the user-specified amount when the
	// order was created.
	//
	// This is an "in-memory" volume that will be reduced when matching the order.
	// When the volume is reduced to 0, the order is considered fulfilled.
	Volume uint64
}

func New(ID string, pair string, price uint64, volume uint64, side OrderSide) (*Order, error) {
	if price <= 0 {
		return nil, fmt.Errorf("Order price must be positive")
	}

	if volume <= 0 {
		return nil, fmt.Errorf("Order volume must be positive")
	}

	if pair == "" {
		return nil, fmt.Errorf("Order pair must not be empty")
	}

	if ID == "" {
		return nil, fmt.Errorf("Order ID must not be empty")
	}

	return &Order{
		ID:     ID,
		Price:  price,
		Volume: volume,
		Side:   side,
	}, nil
}

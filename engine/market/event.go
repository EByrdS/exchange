package market

import (
	"exchange/engine/order"
	"time"
)

// A VolumeEvent signals a change in the volume of a particular price
// in a market.
type VolumeEvent struct {
	// The market pair name. e.g. "USB/GBP"
	Pair string

	// The side of the price, whether it was a sell or buy change.
	Side order.OrderSide

	// The price the volume is referring to.
	Price uint64

	// The actual volume, NOT the volume delta.
	// Prices that no longer have a volume will contain 0 in this field.
	Volume uint64

	// The time of the event.
	Timestamp time.Time
}

// MatchEvent is an event that captures which orders were matched by the engine.
type MatchEvent struct {
	// The market pair name.
	Pair string

	// The ID of the taker order.
	TakerOrderID string

	// The match type of the taker order.
	TakerMatchType order.MatchType

	// The ID of the maker order.
	MakerOrderID string

	// The match type of the maker order.
	MakerMatchType order.MatchType

	// The volume matched for this particular taker-maker match.
	MatchedVolume uint64

	// The settlement price, given by the maker order
	SettlementPrice uint64

	// The time of the event.
	Timestamp time.Time
}

// The type of an order event.
type OrderEventType int

const (
	// An order was accepted by the engine and successfully inserted in the order
	// book. Taker orders are not inserted, only fulfilled partially or totally.
	MakerOrderInserted OrderEventType = iota

	// An order was cancelled, normally by an explicit user action.
	OrderCancelled

	// An order failed validation and it was rejected, so it did not enter the market.
	OrderRejected

	// An taker order could not be fulfilled due to market insolvency.
	TakerOrderUnfulfilled
)

// OrderEvent signals events related to order movements.
type OrderEvent struct {
	// The type of the order event.
	Type OrderEventType

	// The ID of the corresponding order. For events that involve multiple orders,
	// multiple OrderEvents will be fired so that each order can have its own
	// timeline.
	OrderID string

	// The time of the event.
	Timestamp time.Time
}

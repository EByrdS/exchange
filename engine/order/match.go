package order

// The type of a match
type MatchType int

const (
	// An order was fulfilled completely, regardless if this is the only match or
	// the last of multiple partial matches.
	//
	// Each order must have a single OrderFulfilled event.
	//
	// E.g.:
	//
	//  [fulfilled]
	//
	// or:
	//
	//  [partial, partial, partial, fulfilled]
	OrderFulfilled MatchType = iota

	// An order was partially fulfilled but was not completed.
	//
	// Use OrderFulfilled for the last partial fulfillment that completes an order.
	//
	// E.g.:
	//
	//  [partial, partial, partial, fulfilled]
	OrderPartiallyFulfilled
)

// Match is the result of matching one order, containing the volume matched
// and whether the match was in full or not.
type Match struct {
	// The maker order matched
	MakerOrder *Order

	// The type of the match
	Type MatchType

	// The volume taken from the maker order
	VolumeTaken uint64
}

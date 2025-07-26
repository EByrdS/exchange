package orderbook_test

type volumeCallbackParams struct {
	Price  uint64
	Volume uint64
}

type volumeCallbackTracker struct {
	history []*volumeCallbackParams
}

func (t *volumeCallbackTracker) reset() {
	t.history = []*volumeCallbackParams{}
}

func (t *volumeCallbackTracker) call(price uint64, volume uint64) {
	t.history = append(t.history, &volumeCallbackParams{price, volume})
}

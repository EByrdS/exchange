package engineserver

type MarketSymbol struct {
	Base  string
	Trade string
}

func (m *MarketSymbol) Name() string {
	return m.Base + "/" + m.Trade
}

func (m *MarketSymbol) Topic() string {
	return "engine." + m.Base + "." + m.Trade
}

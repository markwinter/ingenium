package ingenium

type OrderType string
type Side string
type Signal string

const (
	MARKET          OrderType = "MARKET"
	BUY             Side      = "BUY"
	SELL            Side      = "SELL"
	LONG            Signal    = "LONG"
	SHORT           Signal    = "SHORT"
	DataEventType             = "ingenium.ingestor.data"
	SignalEventType           = "ingenium.strategy.signal"
	OrderEventType            = "ingenium.portfolio.order"
)

type SignalEvent struct {
	Symbol string
	Signal string
}

type DataEvent struct {
	Symbol     string
	Period     string
	OpenPrice  string
	ClosePrice string
	MaxPrice   string
	MinPrice   string
	Volume     string
}

type OrderEvent struct {
	Side     Side
	Quantity string
	Symbol   string
	Type     OrderType
}

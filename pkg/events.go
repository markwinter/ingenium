package ingenium

type OrderType string
type Side string
type Signal string
type DataType string

const (
	MARKET          OrderType = "MARKET"
	BUY             Side      = "BUY"
	SELL            Side      = "SELL"
	LONG            Signal    = "LONG"
	SHORT           Signal    = "SHORT"
	DataEventType             = "ingenium.ingestor.data"
	SignalEventType           = "ingenium.strategy.signal"
	OrderEventType            = "ingenium.portfolio.order"

	DataTypeOhlc = "data.type.ohlc"
)

type SignalEvent struct {
	Symbol string
	Signal string
}

type DataEvent struct {
	Type      DataType
	Symbol    string
	Timestamp string
	Data      any
}

type DataOhlc struct {
	Open   string
	High   string
	Low    string
	Close  string
	Volume string
	Period string // Period defines the time range of this data e.g. 5 min candle
}

type OrderEvent struct {
	Side     Side
	Quantity string
	Symbol   string
	Type     OrderType
}

type NewEventType struct {
	Key string
}

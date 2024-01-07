package ingenium

import "time"

type OrderType string
type Side string
type Signal string
type DataType string
type TimeInForce string

const (
	MARKET             OrderType   = "MARKET"
	BUY                Side        = "BUY"
	SELL               Side        = "SELL"
	DAY                TimeInForce = "DAY"
	GTC                TimeInForce = "GTC"
	SignalLong         Signal      = "ingenium.signal.long"
	SignalShort        Signal      = "ingenium.signal.short"
	DataEventType                  = "ingenium.ingestor.data"
	SignalEventType                = "ingenium.strategy.signal"
	OrderEventType                 = "ingenium.portfolio.order"
	ExecutionEventType             = "ingenium.executor.execution"

	DataTypeOhlc = "data.type.ohlc"
)

type SignalEvent struct {
	Symbol    string
	Signal    Signal
	Timestamp time.Time
}

type DataEvent struct {
	Type      DataType
	Symbol    string
	Timestamp time.Time
	Ohlc      DataOhlc `json:",omitempty"`
}

type DataOhlc struct {
	Open      string
	High      string
	Low       string
	Close     string
	Volume    string
	Period    string // Period defines the time range of this data e.g. 5 min candle
	Timestamp time.Time
}

type OrderEvent struct {
	Id        string
	Timestamp time.Time

	Side        Side
	Quantity    string
	Symbol      string
	Type        OrderType
	TimeInForce TimeInForce
}

type ExecutionEvent struct {
	OrderId   string
	Quantity  string
	Symbol    string
	Timestamp time.Time
}

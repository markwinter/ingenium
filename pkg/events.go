package ingenium

import "time"

type OrderType string
type Side string
type Signal string
type DataType string
type TimeInForce string

const (
	MarketOrderType OrderType   = "market"
	LimitOrderType  OrderType   = "limit"
	BuySide         Side        = "buy"
	SellSide        Side        = "sell"
	DayTimeInForce  TimeInForce = "day"
	GtcTimeInForce  TimeInForce = "gtc"

	SignalLong  Signal = "ingenium.signal.long"
	SignalShort Signal = "ingenium.signal.short"

	DataEventType      = "ingenium.ingestor.data"
	SignalEventType    = "ingenium.strategy.signal"
	OrderEventType     = "ingenium.portfolio.order"
	ExecutionEventType = "ingenium.executor.execution"

	DataTypeOhlc = "data.type.ohlc"
)

type EventMetadata struct {
	Id        string
	Timestamp time.Time
}

type SignalEvent struct {
	EventMetadata

	Symbol string
	Signal Signal
}

type DataEvent struct {
	EventMetadata

	Type   DataType
	Symbol string
	Ohlc   DataOhlc `json:",omitempty"`
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
	EventMetadata

	Symbol   string
	Side     Side
	Quantity string
	Type     OrderType

	TimeInForce   TimeInForce
	ExtendedHours bool

	LimitPrice   string
	StopPrice    string
	TakeProfit   string
	TrailPrice   string
	TrailPercent string
	StopLoss     StopLoss
}

type StopLoss struct {
	LimitPrice string
	StopPrice  string
}

type ExecutionEvent struct {
	EventMetadata

	ExecutionTimestamp time.Time
	OrderId            string
	Symbol             string
	Quantity           string
	Price              string
}

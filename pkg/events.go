package ingenium

import "time"

type OrderType string
type Side string
type Signal string
type DataType string
type TimeInForce string
type TradeUpdate string

const (
	TimestampFormat = "2006-01-02T15:00:00"

	MarketOrderType OrderType   = "market"
	LimitOrderType  OrderType   = "limit"
	BuySide         Side        = "buy"
	SellSide        Side        = "sell"
	DayTimeInForce  TimeInForce = "day"
	GtcTimeInForce  TimeInForce = "gtc"

	TradeFilled      TradeUpdate = "filled"
	TradePartialFill TradeUpdate = "partial_fill"
	TradeCanceled    TradeUpdate = "canceled"
	TradeRejected    TradeUpdate = "rejected"
	TradeAccepted    TradeUpdate = "accepted"
	TradeNew         TradeUpdate = "new"
	TradeExpired     TradeUpdate = "expired"
	TradeReplaced    TradeUpdate = "replaced"
	TradeHeld        TradeUpdate = "held"

	SignalLong  Signal = "ingenium.signal.long"
	SignalShort Signal = "ingenium.signal.short"

	DataEventType      = "ingenium.ingestor.data"
	SignalEventType    = "ingenium.strategy.signal"
	OrderEventType     = "ingenium.portfolio.order"
	ExecutionEventType = "ingenium.executor.execution"

	DataTypeOhlc = "data.type.ohlc"
)

// EventMetadata gets filled in for you when using the component clients' Send methods
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

	Update TradeUpdate
}

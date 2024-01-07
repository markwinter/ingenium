# Ingenium

### Prototyping Stage

This project is still in prototype stage.

---

Ingenium is a trading system built on top of NATS. It provides a common library (`go get github.com/markwinter/ingenium/pkg`) that allows quick creation of new components. It's expected that end users create their own components, e.g. ingestors and strategies, using the common library that handles correct typing etc. so that components can easily communicate with each other.

Ingenium is event-driven using [NATS](https://nats.io/)

Ingenium will come with telemetry built in using OpenTelemetry.

There will also be a web component in the future to view and manage the current state of the system such as
- View currently running components e.g. enabled strategies, data ingestors
- Launch backtests
- View portfolio stats and transaction history

## Components

A simplified diagram of the system is below. In reality you can have multiple of each component running at the same time.


           1.                                         2.
          ┌────────────┐                   ┌────────────┐
          │            │                   │            │
          │  Ingestor  │       Data Event  │  Strategy  │
          │            │     ┌─────────────►            │
          └──────┬─────┘     │             └───────────┬┘
                 │           │                         │
                 │           │                         │
                 │         ┌─┴──────────────┐          │
                 │         │                │          │
                 └─────────►  Event Broker  ◄──────────┘
               Data Event  │                │     Signal Event
                           └▲─▲────────────┬┘
                            │ │            │
                            │ │            └─────────┐
                    ┌───────┘ └────────────┐         │Signal Event
                    │     Order Event /    │         │
    ┌───────────────▼───┐ Execution Event ┌▼─────────▼──┐
    │                   │                 │             │
    │  Order Executor   │                 │  Portfolio  │
    │                   │                 │             │
    └───────────────────┘                 └─────────────┘
     4.                                               3.


### Ingestors

Ingestors feed market data into the system. The component produces a data event `ingenium.DataEvent` for each market data
and sends it to a subject in the format: `ingenium.ingestor.data.<stock-symbol>`

Examples of ingestors:

- One-shot jobs that read historical data from a CSV file or an API
- A long running binary that reads from a live market exchange

### Strategies

Strategies receive data events from Ingestors by subscribing to a data event subject `ingenium.ingestor.data.<stock-symbol>` (or wildcard data `ingenium.ingestor.data.*`). Strategies produce signal events `ingenium.SignalEvent` based on an implemented trading strategy and send it to subject `ingenium.strategy.signal`

### Portfolios

Portfolios receive signal events from Strategies and decide based on several factors whether to generate a market order event `ingenium.OrderEvent` and send it to the subject
`ingenium.portfolio.order`.

### Order Executors

Order Executors receive market order events from Portfolios and execute the appropriate order
on the exchange. They also return order execution events back to the Portfolio using the subject `ingenium.executor.execution`. In the case of partial fills, multiple messages are sent containing the original order id so the Portfolio can track its position.

## Events

All events are defined as Golang structs. Currently they are serialized to JSON.

Below is a list of all events in the system and their spec

### Market Data

Type: `ingenium.ingestor.data`

```GO
type DataEvent struct {
  Type      DataType
  Symbol    string
  Timestamp string
  Ohlc      DataOhlc `json:",omitempty"`
}

type DataOhlc struct {
  Open   string
  High   string
  Low    string
  Close  string
  Volume string
}
```

### Signal

Type: `ingenium.strategy.signal`

```GO
type Signal struct {
  Symbol string
  Signal string
}
```

### Market Order

Type: `ingenium.portfolio.order`

```GO
type Order struct {
	Side     Side
	Quantity big.Decimal
	Symbol   string
	Type     OrderType
}

type OrderType string
type Side string

const (
	MARKET OrderType = "MARKET"
	BUY    Side      = "BUY"
	SELL   Side      = "SELL"
)
```

### Order Execution

Type: `ingenium.executor.execution`

```GO
```


## Roadmap

- Finish building examples of each component
- Add component generators/functions to `/pkg`
- Backtesting lib
- A database component to record events, trades, portfolio data etc.
- Allow deploys to local, or to Kube cluster (similar to [Service Weaver](https://serviceweaver.dev/))

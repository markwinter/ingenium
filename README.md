# Ingenium

**-- Still heavily in prototyping stage --**

Ingenium is a cloud native electronic trading system built on top of Kubernetes and Knative Eventing.

Ingenium is event-driven, using CloudEvents and Knative Eventing to pass data between components.

Ingenium comes with telemetry built in using OpenTelemetry.

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
          └──────┬─────┘     │             └─────┬──────┘
                 │           │                   │
                 │           │                   │
                 │         ┌─┴────────┐          │
                 │         │          │          │
                 └─────────►  Broker  ◄──────────┘
               Data Event  │          │     Signal Event
                           └─────────┬┘
                                     │
                                     └───────────────┐
                                                     │Signal Event
                                                     │
    ┌───────────────────┐ Execution Event ┌──────────▼──┐
    │                   ├─────────────────►             │
    │  Order Executor   │                 │  Portfolio  │
    │                   ◄─────────────────┤             │
    └───────────────────┘   Order Event   └─────────────┘
     4.                                               3.


### Broker

The broker handles receiving and sending events between components.

By default Kafka is used for the broker.
This can be changed to something else like GCP PubSub based on your needs.
Keep in mind that message ordering is generally required for strategies to work correctly.

If you have multiple of the same component e.g. Strategies, the same event will be delivered to each component.
When you create a `Trigger`, Kafka will send any existing messages retained in the topic.

### Ingestors

Ingestors feed market data into the system. The component produces a data event for each market data
which gets sent to the Broker.

Ingestors can be either:

- One-Shot Kubernetes Jobs that read from files
- A long running Kubernetes Deployment that reads from a real market exchange

### Strategies

Strategies receive market data events from Ingestors and produce signal events based on an implemented
trading strategy.

### Portfolios

Portfolios receive signal events from Strategies and decide based on several factors such as
remaining balance, risk assessment etc. whether to generate a market order event. When generating a market
order appropriate order sizing also takes place. Portfolios also manage open positions.

### Order Executors

Order Executors receive market order events from Portfolios and execute the appropriate order
on the exchange. They also return order execution events back to the Portfolio.

## Events

All events are CloudEvents generated using the CloudEvents SDKs. Currently they all serialized to JSON.

Below is a list of all Events in the system and their spec.

### Market Data

```GO
type DataEvent struct {
  Symbol     string
	Period     string
	OpenPrice  string
	ClosePrice string
	MaxPrice   string
	MinPrice   string
	Volume     string
}
```

### Signal

```GO
type Signal struct {
  Symbol string
	Signal string
}
```

### Market Order

### Order Execution

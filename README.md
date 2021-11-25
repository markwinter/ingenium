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

### Broker

By default and for testing the broker is Knative Eventing's Multi-Tenant Channel Broker using in-memory channels.
For production systems this can be easily changed to something more suitable like Kafka or GCP PubSub.

### Ingestors

Ingestors feed market data into the system. The component produces a CloudEvent for each market data
which gets sent to the Broker.

Ingestors can be either:
    - One-Shot Kubernetes Jobs that read from files
    - A long running Kubernetes Deployment that reads from a real market exchange

### Strategies

Strategies receive market data events from Ingestors (via the Broker) and produce signals based on an implemented
trading strategy. This signal is created as a CloudEvent which gets set to the Broker.

### Portfolios

Portfolios receive signal events from Strategies (via the Broker) and decide based on several factors such as
remaining balance, risk assessment, outstanding orders, etc. whether to generate a market order event. When generating a market
order appropriate order sizing also takes place. Portfolios also manage open positions.

### Order Executors

Order Executors receive market order events from Portfolios (via the Broker) and execute the appropriate order
on the exchange. They also return order execution events back to the Portfolio.

## Events

All events are CloudEvents generated using the CloudEvents SDKs. Below is a list of all Events in the system and their spec.

### Market Data

### Signal

### Market Order

### Order Execution

package main

import (
	"context"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/sdcoffey/techan"
)

type RsiStrategy struct {
	series    *techan.TimeSeries
	strategy  techan.Strategy
	indicator techan.Indicator
	record    *techan.TradingRecord
	index     int
}

func MakeRsiStrategy() *RsiStrategy {
	series := techan.NewTimeSeries()
	closePrices := techan.NewClosePriceIndicator(series)
	rsi := techan.NewRelativeStrengthIndexIndicator(closePrices, 14)

	// Enter position when RSI goes below 30
	entryRule := techan.And(
		techan.NewCrossDownIndicatorRule(rsi, techan.NewConstantIndicator(30)),
		techan.PositionNewRule{},
	)

	// Exit position when RSI goes above 70
	exitRule := techan.And(
		techan.NewCrossUpIndicatorRule(techan.NewConstantIndicator(70), rsi),
		techan.PositionOpenRule{},
	)

	strategy := techan.RuleStrategy{
		UnstablePeriod: 14,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}

	return &RsiStrategy{
		series:    series,
		strategy:  strategy,
		indicator: rsi,
		record:    techan.NewTradingRecord(),
		index:     0,
	}
}

var rsiStrategy *RsiStrategy

func receive(event cloudevents.Event) {
	log.Printf("%s", event)
}

func main() {
	rsiStrategy = MakeRsiStrategy()

	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Print("Starting HTTP Receiver")

	log.Fatal(c.StartReceiver(context.Background(), receive))
}

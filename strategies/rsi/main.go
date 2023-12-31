package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	ingenium "github.com/markwinter/ingenium/pkg"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

type RsiStrategy struct {
	series    *techan.TimeSeries
	strategy  techan.Strategy
	indicator techan.Indicator
	record    *techan.TradingRecord
	index     int
	client    cloudevents.Client
}

func MakeRsiStrategy() RsiStrategy {
	series := techan.NewTimeSeries()
	closePrices := techan.NewClosePriceIndicator(series)
	rsi := techan.NewRelativeStrengthIndexIndicator(closePrices, 14)

	// Enter position when RSI goes below 30
	entryRule := techan.NewCrossDownIndicatorRule(rsi, techan.NewConstantIndicator(30))

	// Exit position when RSI goes above 70
	exitRule := techan.NewCrossUpIndicatorRule(techan.NewConstantIndicator(70), rsi)

	strategy := techan.RuleStrategy{
		UnstablePeriod: 14,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}

	client, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	return RsiStrategy{
		series:    series,
		strategy:  strategy,
		indicator: rsi,
		record:    techan.NewTradingRecord(),
		client:    client,
	}
}

var rsiStrategy RsiStrategy = MakeRsiStrategy()
var broker string

func sendEvent(symbol, signal string) {
	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetTime(time.Now())
	event.SetSource(fmt.Sprintf("ingenium/strategy/rsi/%s", os.Getenv("HOSTNAME")))
	event.SetType(ingenium.SignalEventType)

	event.SetData(cloudevents.ApplicationJSON, ingenium.SignalEvent{Signal: signal, Symbol: symbol})

	ctx := cloudevents.ContextWithTarget(context.Background(), broker)
	if result := rsiStrategy.client.Send(ctx, event); cloudevents.IsUndelivered(result) {
		log.Printf("failed to send, %v", result)
	}
}

func receive(event cloudevents.Event) {
	var dataEvent ingenium.DataEvent
	if err := json.Unmarshal(event.Data(), &dataEvent); err != nil {
		log.Printf("[%s] Failed to unmarshal event: %v", event.ID(), err)
		return
	}

	date, err := time.Parse("2006-01-02", dataEvent.Period)
	if err != nil {
		log.Printf("[%s] Failed to parse date: %v", event.ID(), err)
		return
	}
	period := techan.NewTimePeriod(date, time.Hour*24)

	candle := techan.NewCandle(period)
	candle.OpenPrice = big.NewFromString(dataEvent.OpenPrice)
	candle.ClosePrice = big.NewFromString(dataEvent.ClosePrice)
	candle.MaxPrice = big.NewFromString(dataEvent.MaxPrice)
	candle.MinPrice = big.NewFromString(dataEvent.MinPrice)
	candle.Volume = big.NewFromString(dataEvent.Volume)

	rsiStrategy.series.AddCandle(candle)

	if rsiStrategy.strategy.ShouldEnter(rsiStrategy.index, rsiStrategy.record) {
		log.Printf("[RSI] %s | LONG | RSI: %v", dataEvent.Symbol, rsiStrategy.indicator.Calculate(rsiStrategy.index))
		sendEvent(dataEvent.Symbol, "long")
	} else if rsiStrategy.strategy.ShouldExit(rsiStrategy.index, rsiStrategy.record) {
		log.Printf("[RSI] %s | SHORT | RSI: %v", dataEvent.Symbol, rsiStrategy.indicator.Calculate(rsiStrategy.index))
		sendEvent(dataEvent.Symbol, "short")
	}

	rsiStrategy.index++
}

func main() {
	b := flag.String("broker", "", "URL of broker to send events to")

	flag.Parse()

	if *b == "" {
		log.Fatalf("broker url was not given")
	}
	broker = *b

	client, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Print("Starting HTTP Receiver")

	log.Fatal(client.StartReceiver(context.Background(), receive))
}

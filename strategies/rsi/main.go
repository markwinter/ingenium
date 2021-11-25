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
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

type RsiStrategy struct {
	series    *techan.TimeSeries
	strategy  techan.Strategy
	indicator techan.Indicator
	record    *techan.TradingRecord
	index     int
}

type DataEvent struct {
	Period     string
	OpenPrice  string
	ClosePrice string
	MaxPrice   string
	MinPrice   string
	Volume     string
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
var client cloudevents.Client
var broker string

func sendEvent(signal string) {
	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetTime(time.Now())
	event.SetSource(fmt.Sprintf("ingenium/strategy/rsi/%s", os.Getenv("HOSTNAME")))
	event.SetType("ingenium.strategy.signal")

	event.SetData(cloudevents.ApplicationJSON, map[string]string{"signal": signal})

	ctx := cloudevents.ContextWithTarget(context.Background(), broker)
	if result := client.Send(ctx, event); cloudevents.IsUndelivered(result) {
		log.Printf("failed to send, %v", result)
	}
}

func receive(event cloudevents.Event) {
	var dataEvent DataEvent
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
		log.Printf("[RSI] LONG | RSI: %v", rsiStrategy.indicator.Calculate(rsiStrategy.index))
		sendEvent("long")
	} else if rsiStrategy.strategy.ShouldExit(rsiStrategy.index, rsiStrategy.record) {
		log.Printf("[RSI] SHORT | RSI: %v", rsiStrategy.indicator.Calculate(rsiStrategy.index))
		sendEvent("short")
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

	rsiStrategy = MakeRsiStrategy()

	var err error
	client, err = cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Print("Starting HTTP Receiver")

	log.Fatal(client.StartReceiver(context.Background(), receive))
}

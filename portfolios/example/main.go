package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	ingenium "github.com/markwinter/ingenium/pkg"
	"github.com/segmentio/ksuid"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

type Portfolio struct {
	currency  big.Decimal
	positions map[string]*techan.TradingRecord
	client    cloudevents.Client
}

func MakePortfolio(m float64) Portfolio {
	client, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	return Portfolio{
		currency: big.NewDecimal(m),
		client:   client,
	}
}

var portfolio Portfolio
var broker string

func getPosition(symbol string) *techan.TradingRecord {
	if val, ok := portfolio.positions[symbol]; ok {
		return val
	}

	record := techan.NewTradingRecord()
	portfolio.positions[symbol] = record

	return record
}

func generateOrder(side ingenium.Side, symbol string, quantity big.Decimal) ingenium.OrderEvent {
	return ingenium.OrderEvent{
		Side:     side,
		Symbol:   symbol,
		Quantity: quantity.String(),
		Type:     ingenium.MARKET,
	}
}

func sendOrder(order ingenium.OrderEvent) {
	event := cloudevents.NewEvent()
	event.SetID(fmt.Sprintf("order_%s", ksuid.New()))
	event.SetTime(time.Now())
	event.SetSource(fmt.Sprintf("ingenium/portfolio/example/%s", os.Getenv("HOSTNAME")))
	event.SetType("ingenium.portfolio.order")

	if err := event.SetData(cloudevents.ApplicationJSON, order); err != nil {
		log.Printf("failed to set data on event: %v", err)
		return
	}

	ctx := cloudevents.ContextWithTarget(context.Background(), broker)
	if result := portfolio.client.Send(ctx, event); cloudevents.IsUndelivered(result) {
		log.Printf("failed to send, %v", result)
	}
}

func long(symbol string) {
	position := getPosition(symbol)

	// Example portfolio doesn't increase position after initial position
	if position.CurrentPosition().IsOpen() {
		return
	}

	quantity := big.NewDecimal(1.0)

	order := generateOrder(ingenium.BUY, symbol, quantity)
	sendOrder(order)
}

func short(symbol string) {
	position := getPosition(symbol)

	// Example portfolio does not allow margin so only sell open long positions
	if !position.CurrentPosition().IsOpen() {
		return
	}

	// Example portfolio always closes position completely
	quantity := position.CurrentPosition().EntranceOrder().Amount

	order := generateOrder(ingenium.SELL, symbol, quantity)
	sendOrder(order)
}

func handleSignal(event cloudevents.Event) {
	var signalEvent ingenium.SignalEvent
	if err := json.Unmarshal(event.Data(), &signalEvent); err != nil {
		log.Printf("[%s] Failed to unmarshal event: %v", event.ID(), err)
		return
	}

	if signalEvent.Signal == ingenium.SignalLong {
		long(signalEvent.Symbol)
	} else if signalEvent.Signal == ingenium.SignalShort {
		short(signalEvent.Symbol)
	}
}

func handleExecution(event cloudevents.Event) {

}

func receive(event cloudevents.Event) {
	switch event.Type() {
	case "ingenium.strategy.signal":
		handleSignal(event)
	case "ingenium.executor.execution":
		handleExecution(event)
	}
}

func main() {
	b := flag.String("broker", "", "URL of broker to send events to")
	m := flag.Float64("money", 10000.0, "Starting currency")

	flag.Parse()

	if *b == "" {
		log.Fatalf("broker url was not given")
	}
	broker = *b

	portfolio = MakePortfolio(*m)

	listenClient, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Print("Starting HTTP Receiver")

	log.Fatal(listenClient.StartReceiver(context.Background(), receive))
}

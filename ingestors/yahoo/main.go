package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	ingenium "github.com/markwinter/ingenium/pkg"
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

var (
	dataSourceName = fmt.Sprintf("ingenium/ingestor/csv/%s", os.Getenv("HOSTNAME"))
)

func main() {
	broker := flag.String("broker", "", "URL of broker to send events to")
	symbol := flag.String("symbol", "", "Security symbol")
	timeFrameRaw := flag.String("timeframe", "24h", "Timeframe given in Golang time string format e.g. 1h, 1d, 12h")

	flag.Parse()

	if *broker == "" {
		log.Fatalf("broker url was not given")
	}

	if *symbol == "" {
		log.Fatalf("symbol was not given")
	}

	timeFrame, err := time.ParseDuration(*timeFrameRaw)
	if err != nil {
		log.Fatalf("timeframe %q is invalid: %v", *timeFrameRaw, err)
	}

	client, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create CloudEvents client: %v", err)
	}

	log.Printf("Getting OHLC for %q with timeframe %q", *symbol, timeFrame.String())

	params := &chart.Params{
		Symbol:   *symbol,
		Interval: datetime.OneHour,
	}

	iter := chart.Get(params)

	for iter.Next() {
		d := iterBarToDataEvent(*symbol, *timeFrameRaw, iter.Bar())
		fmt.Printf("%v", d)
		sendEvent(client, dataSourceName, *broker, d)
	}

	if err := iter.Err(); err != nil {
		fmt.Println(err)
	}
}

func iterBarToDataEvent(symbol, timeFrame string, bar *finance.ChartBar) ingenium.DataEvent {
	return ingenium.DataEvent{
		Symbol:    symbol,
		Type:      ingenium.DataTypeOhlc,
		Timestamp: fmt.Sprintf("%d", bar.Timestamp),
		Data: ingenium.DataOhlc{
			Open:   bar.Open.String(),
			High:   bar.High.String(),
			Low:    bar.Low.String(),
			Close:  bar.Close.String(),
			Volume: fmt.Sprintf("%d", bar.Volume),
			Period: timeFrame,
		},
	}
}

func sendEvent(client cloudevents.Client, sourceName, broker string, data ingenium.DataEvent) {
	event, err := ingenium.ConvertDataEvent(data, sourceName)
	if err != nil {
		log.Printf("failed to convert to cloudevent: %v", err)
		return
	}

	ctx := cloudevents.ContextWithTarget(context.Background(), broker)
	if result := client.Send(ctx, event); cloudevents.IsUndelivered(result) {
		log.Printf("failed to send, %v", result)
	}
}

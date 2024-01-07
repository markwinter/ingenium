package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	alpaca "github.com/markwinter/ingenium/examples/ingestors/alpaca-historical"
	ingenium "github.com/markwinter/ingenium/pkg"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	symbol := flag.String("symbol", "", "Security symbol")
	startString := flag.String("start", "", "Start date of of data to retrieve")
	endString := flag.String("end", "", "End date of data to retrieve")
	timeframeString := flag.String("timeframe", "1h", "aggregated timeframe of ohlcv")

	flag.Parse()

	if *symbol == "" {
		log.Fatalf("symbol was not given")
	}

	end, err := parseDateArg(*endString, time.Now())
	if err != nil {
		log.Fatalf("failed to parse start argument: %v", err)
	}

	start, err := parseDateArg(*startString, end.Add(-12*time.Hour))
	if err != nil {
		log.Fatalf("failed to parse start argument: %v", err)
	}

	fmt.Printf("start: %s\n", start)
	fmt.Printf("end: %s\n", end)
	fmt.Printf("timeframe: %s\n", *timeframeString)

	ingestor := alpaca.NewAlpacaHistoricalIngestor(*symbol, start, end, *timeframeString)
	defer ingestor.Cleanup()

	ingestor.IngestData()
}

func parseDateArg(value string, defaultValue time.Time) (time.Time, error) {
	if value == "" {
		return defaultValue, nil
	}

	return time.Parse(ingenium.TimestampFormat, value)
}

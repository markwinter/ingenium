package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	alpaca "github.com/markwinter/ingenium/examples/ingestors/alpaca-historical"
	exampleportfolio "github.com/markwinter/ingenium/examples/portfolios/example"
	rsi "github.com/markwinter/ingenium/examples/strategies/rsi"
	"github.com/markwinter/ingenium/pkg/backtest"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// prints all events ingenium.> to see in the console what happens during a backtest
	_ = MakeEventPrinter()

	dataStart := time.Date(2024, 01, 04, 04, 00, 00, 00, time.UTC)
	dataEnd := time.Date(2024, 01, 04, 23, 00, 00, 00, time.UTC)

	symbol := "CPNG"

	backtest := backtest.NewBacktest(
		// Run locally or deploy to kubernetes
		backtest.WithDeploymentType(backtest.DeploymentLocal),
		backtest.WithIngestor(alpaca.NewAlpacaHistoricalIngestor(symbol, dataStart, dataEnd, "1h")),
		backtest.WithStrategy(rsi.NewRsiStrategy(symbol)),
		backtest.WithPortfolio(exampleportfolio.NewPortfolio(1000)),
	)
	defer backtest.Cleanup()

	backtest.Run()

	// Get backtest statistics e.g. orders, gain/loss, max drawdown etc.

	// For now wait for Ctrl+C because NATs cleanup isnt done correctly yet and need to wait for all messages to be processed

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done
}

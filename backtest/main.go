package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	alpaca "github.com/markwinter/ingenium/ingestors/alpacahistorical/pkg"
	ingenium "github.com/markwinter/ingenium/pkg"
	rsi "github.com/markwinter/ingenium/strategies/rsi/pkg"
)

type DeploymentType string

const (
	DeploymentLocal DeploymentType = "local"
	// DeploymentKube DeploymentType = "kube"
)

type Backtest struct {
	Deployment DeploymentType
	Ingestors  []ingenium.Ingestor
	Strategies []ingenium.Strategy
	Portfolio  ingenium.Portfolio
	Executor   ingenium.Executor

	printer EventPrinter
}

type BacktestOption func(*Backtest)

func WithDeploymentType(t DeploymentType) BacktestOption {
	return func(b *Backtest) {
		b.Deployment = t
	}
}

func WithIngestor(ingestor ingenium.Ingestor) BacktestOption {
	return func(b *Backtest) {
		b.Ingestors = append(b.Ingestors, ingestor)
	}
}

func WithStrategy(strategy ingenium.Strategy) BacktestOption {
	return func(b *Backtest) {
		b.Strategies = append(b.Strategies, strategy)
	}
}

func WithPortfolio(portfolio ingenium.Portfolio) BacktestOption {
	return func(b *Backtest) {
		b.Portfolio = portfolio
	}
}

func WithExecutor(executor ingenium.Portfolio) BacktestOption {
	return func(b *Backtest) {
		b.Executor = executor
	}
}

func NewBacktest(options ...BacktestOption) *Backtest {
	b := &Backtest{
		Deployment: DeploymentLocal,
		printer:    MakeEventPrinter(),
	}

	for _, opt := range options {
		opt(b)
	}

	return b
}

func (b *Backtest) Run() {
	time.Sleep(2 * time.Second)

	for _, ingestor := range b.Ingestors {
		ingestor.IngestData()
	}
}

func (b *Backtest) Cleanup() {
	for _, ingestor := range b.Ingestors {
		ingestor.Cleanup()
	}

	for _, strategy := range b.Strategies {
		strategy.Cleanup()
	}

	//b.Portfolio.Cleanup()
	//b.Executor.Cleanup()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dataStart := time.Date(2024, 01, 04, 04, 00, 00, 00, time.UTC)
	dataEnd := time.Date(2024, 01, 04, 23, 00, 00, 00, time.UTC)

	symbol := "CPNG"

	backtest := NewBacktest(
		// Run locally or deploy to kubernetes
		WithDeploymentType(DeploymentLocal),
		WithIngestor(alpaca.NewAlpacaHistoricalIngestor(symbol, dataStart, dataEnd, "1h")),
		WithStrategy(rsi.NewRsiStrategy(symbol)),
	)
	defer backtest.Cleanup()

	backtest.Run()

	// Get backtest statistics e.g. orders, gain/loss, max drawdown etc.

	// For now wait for Ctrl+C because NATs cleanup isnt done correctly yet and need to wait for all messages to be processed

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done
}

package backtest

import (
	ingenium "github.com/markwinter/ingenium/pkg"
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
	}

	for _, opt := range options {
		opt(b)
	}

	return b
}

func (b *Backtest) Run() {
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

	if b.Portfolio != nil {
		b.Portfolio.Cleanup()
	}

	if b.Executor != nil {
		b.Executor.Cleanup()
	}
}

package rsi

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/cinar/indicator"
	ingenium "github.com/markwinter/ingenium/pkg"
	"github.com/markwinter/ingenium/pkg/strategy"
)

type RsiStrategy struct {
	strategyClient *strategy.StrategyClient

	minWindow int
	buyAt     float64
	sellAt    float64

	closingPrices map[string][]float64
}

func NewRsiStrategy(symbols ...string) *RsiStrategy {
	strategy := &RsiStrategy{
		strategyClient: strategy.NewStrategy(),
		minWindow:      2,
		buyAt:          5,
		sellAt:         95,
		closingPrices:  make(map[string][]float64),
	}

	for _, symbol := range symbols {
		if err := strategy.strategyClient.SubscribeToSymbol(symbol, strategy.Receive); err != nil {
			log.Printf("failed to subscribe to symbol data: %v", err)
		}
	}

	return strategy
}

func (s *RsiStrategy) Receive(dataEvent *ingenium.DataEvent) {
	// Check the dataEvent.Type if you are uncertain about the data type
	symbol := dataEvent.Symbol
	data := dataEvent.Ohlc

	close, _ := strconv.ParseFloat(data.Close, 64)

	s.closingPrices[symbol] = append(s.closingPrices[symbol], close)

	if len(s.closingPrices[symbol]) <= s.minWindow {
		return
	}

	_, rsi := indicator.Rsi2(s.closingPrices[symbol])

	if rsi[len(rsi)-1] < s.buyAt {
		fmt.Printf("[RSI] %s | LONG | %v\n", symbol, rsi[len(rsi)-1])

		event := ingenium.SignalEvent{
			Symbol:    dataEvent.Symbol,
			Signal:    ingenium.SignalLong,
			Timestamp: time.Now(),
		}

		if err := s.strategyClient.SendSignalEvent(event); err != nil {
			log.Printf("failed sending signal: %v", err)
		}
	} else if rsi[len(rsi)-1] > s.sellAt {
		fmt.Printf("[RSI] %s | SHORT | %v\n", symbol, rsi[len(rsi)-1])

		event := ingenium.SignalEvent{
			Symbol:    dataEvent.Symbol,
			Signal:    ingenium.SignalShort,
			Timestamp: time.Now(),
		}

		if err := s.strategyClient.SendSignalEvent(event); err != nil {
			log.Printf("failed sending signal: %v", err)
		}
	}
}

func (s *RsiStrategy) Run() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done
}

func (s *RsiStrategy) Cleanup() {
	s.strategyClient.Close()
}

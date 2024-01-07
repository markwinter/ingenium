package rsi

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/cinar/indicator"
	ingenium "github.com/markwinter/ingenium/pkg"
	"github.com/markwinter/ingenium/pkg/strategy"
)

type RsiStrategy struct {
	*strategy.StrategyClient

	minWindow int
	buyAt     float64
	sellAt    float64

	closingPrices map[string][]float64
}

func NewRsiStrategy(symbols ...string) *RsiStrategy {
	s := &RsiStrategy{
		minWindow:     2,
		buyAt:         5,
		sellAt:        95,
		closingPrices: make(map[string][]float64),
	}

	client := strategy.NewStrategyClient(s, symbols)
	s.StrategyClient = client

	return s
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
		//fmt.Printf("[RSI] %s | LONG | %v\n", symbol, rsi[len(rsi)-1])

		event := ingenium.SignalEvent{
			Symbol: dataEvent.Symbol,
			Signal: ingenium.SignalLong,
		}

		if err := s.SendSignalEvent(event); err != nil {
			log.Printf("failed sending signal: %v", err)
		}
	} else if rsi[len(rsi)-1] > s.sellAt {
		//fmt.Printf("[RSI] %s | SHORT | %v\n", symbol, rsi[len(rsi)-1])

		event := ingenium.SignalEvent{
			Symbol: dataEvent.Symbol,
			Signal: ingenium.SignalShort,
		}

		if err := s.SendSignalEvent(event); err != nil {
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
	s.Close()
}

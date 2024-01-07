package exampleportfolio

import (
	"os"
	"os/signal"
	"syscall"

	ingenium "github.com/markwinter/ingenium/pkg"
	"github.com/markwinter/ingenium/pkg/portfolio"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

type ExamplePortfolio struct {
	*portfolio.PortfolioClient

	balance   big.Decimal
	positions map[string]*techan.Position
	record    map[string]*techan.TradingRecord
}

func NewPortfolio(b float64) *ExamplePortfolio {
	p := &ExamplePortfolio{
		balance:   big.NewDecimal(b),
		positions: make(map[string]*techan.Position),
		record:    make(map[string]*techan.TradingRecord),
	}

	p.PortfolioClient = portfolio.NewPortfolioClient(p)

	return p
}

func generateOrder(side ingenium.Side, symbol string, quantity big.Decimal) ingenium.OrderEvent {
	return ingenium.OrderEvent{
		Side:        side,
		Symbol:      symbol,
		Quantity:    quantity.String(),
		Type:        ingenium.MarketOrderType,
		TimeInForce: ingenium.GtcTimeInForce,
	}
}

func (p *ExamplePortfolio) long(symbol string) {
	record, ok := p.record[symbol]
	if !ok {
		record = techan.NewTradingRecord()
	}

	// Example portfolio doesn't increase position after initial position
	if record.CurrentPosition().IsOpen() {
		return
	}

	// Example portfolio always buys just 1 share
	quantity := big.NewDecimal(1.0)

	order := generateOrder(ingenium.BuySide, symbol, quantity)
	p.SendOrder(order)
}

func (p *ExamplePortfolio) short(symbol string) {
	record, ok := p.record[symbol]
	if !ok || !record.CurrentPosition().IsOpen() {
		// Example portfolio does not allow shorting so only sell open long positions
		return
	}

	// Example portfolio always closes position completely
	quantity := record.CurrentPosition().EntranceOrder().Amount

	order := generateOrder(ingenium.SellSide, symbol, quantity)
	p.SendOrder(order)
}

func (p *ExamplePortfolio) ReceiveSignal(event *ingenium.SignalEvent) {
	if event.Signal == ingenium.SignalLong {
		p.long(event.Symbol)
	} else if event.Signal == ingenium.SignalShort {
		p.short(event.Symbol)
	}
}

func (p *ExamplePortfolio) ReceiveExecution(event *ingenium.ExecutionEvent) {
	order := techan.Order{
		//Side:
		Security:      event.Symbol,
		Price:         big.NewFromString(event.Price),
		Amount:        big.NewFromString(event.Quantity),
		ExecutionTime: event.ExecutionTimestamp,
	}

	position, ok := p.positions[event.Symbol]
	if ok {
		order.Side = techan.SELL
		position.Exit(order)
	} else {
		order.Side = techan.BUY
		p.positions[event.Symbol] = techan.NewPosition(order)
	}

	p.record[event.Symbol].Operate(order)
}

func (p *ExamplePortfolio) Run() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done
}

func (p *ExamplePortfolio) Cleanup() {
	p.Close()
}

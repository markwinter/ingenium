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
	positions map[string]*techan.TradingRecord
}

func NewPortfolio(b float64) *ExamplePortfolio {
	p := &ExamplePortfolio{
		balance:   big.NewDecimal(b),
		positions: make(map[string]*techan.TradingRecord),
	}

	p.PortfolioClient = portfolio.NewPortfolioClient(p)

	return p
}

func (p *ExamplePortfolio) getPosition(symbol string) *techan.TradingRecord {
	if val, ok := p.positions[symbol]; ok {
		return val
	}

	record := techan.NewTradingRecord()
	p.positions[symbol] = record

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

func (p *ExamplePortfolio) long(symbol string) {
	position := p.getPosition(symbol)

	// Example portfolio doesn't increase position after initial position
	if position.CurrentPosition().IsOpen() {
		return
	}

	quantity := big.NewDecimal(1.0)

	order := generateOrder(ingenium.BUY, symbol, quantity)
	p.SendOrder(order)
}

func (p *ExamplePortfolio) short(symbol string) {
	position := p.getPosition(symbol)

	// Example portfolio does not allow margin so only sell open long positions
	if !position.CurrentPosition().IsOpen() {
		return
	}

	// Example portfolio always closes position completely
	quantity := position.CurrentPosition().EntranceOrder().Amount

	order := generateOrder(ingenium.SELL, symbol, quantity)
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
}

func (p *ExamplePortfolio) Run() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done
}

func (p *ExamplePortfolio) Cleanup() {
	p.Close()
}

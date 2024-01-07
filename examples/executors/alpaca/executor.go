package alpaca

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	ingenium "github.com/markwinter/ingenium/pkg"
	"github.com/markwinter/ingenium/pkg/executor"
	"github.com/shopspring/decimal"
)

type AlpacaExecutor struct {
	*executor.ExecutorClient

	alpacaClient *alpaca.Client
}

func NewAlpacaExecutor() *AlpacaExecutor {
	e := &AlpacaExecutor{
		alpacaClient: alpaca.NewClient(alpaca.ClientOpts{}),
	}

	e.alpacaClient.StreamTradeUpdatesInBackground(context.Background(), e.tradeHandler)

	e.ExecutorClient = executor.NewExecutorClient(e)

	return e
}

func (e *AlpacaExecutor) ReceiveOrder(order *ingenium.OrderEvent) {
	req := convertIngeniumOrder(order)

	if _, err := e.alpacaClient.PlaceOrder(req); err != nil {
		log.Printf("failed to place order: %v", err)
	}
}

func (e *AlpacaExecutor) tradeHandler(trade alpaca.TradeUpdate) {
	if trade.Event != "fill" {
		return
	}

	// TODO: Handle all trade.Event types

	event := ingenium.ExecutionEvent{
		OrderId:            trade.Order.ClientOrderID,
		Quantity:           trade.Qty.String(),
		Price:              trade.Price.String(),
		ExecutionTimestamp: *trade.Timestamp,
	}

	if err := e.SendExecutionEvent(event); err != nil {
		log.Printf("failed sending execution event: %v", err)
	}
}

func (e *AlpacaExecutor) Cleanup() {
	e.Close()
}

func (e *AlpacaExecutor) Run() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done
}

func convertIngeniumOrder(order *ingenium.OrderEvent) alpaca.PlaceOrderRequest {
	qty, _ := decimal.NewFromString(order.Quantity)

	request := alpaca.PlaceOrderRequest{
		ClientOrderID: order.Id,
		Symbol:        order.Symbol,
		Qty:           &qty,
		Type:          alpaca.OrderType(order.Type),
		Side:          alpaca.Side(order.Side),
		TimeInForce:   alpaca.TimeInForce(order.Type),
	}

	if order.Type == ingenium.LimitOrderType {
		lp, _ := decimal.NewFromString(order.LimitPrice)
		request.LimitPrice = &lp
	}

	return request
}

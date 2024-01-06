package strategy

import (
	"fmt"
	"log"
	"os"
	"strings"

	ingenium "github.com/markwinter/ingenium/pkg"
	"github.com/nats-io/nats.go"
)

type StrategyClient struct {
	natsServer string
	nc         *nats.Conn
	ec         *nats.EncodedConn
}

type StrategyOption func(*StrategyClient)

// WithNatsServer sets the nats-server address. Default is http://0.0.0.0:4022
func WithNatsServer(server string) StrategyOption {
	return func(s *StrategyClient) {
		s.natsServer = server
	}
}

func NewStrategy(opts ...StrategyOption) *StrategyClient {
	var natsServer string

	natsServer = os.Getenv("NATS_SERVER")
	if natsServer == "" {
		natsServer = nats.DefaultURL
	}

	i := &StrategyClient{
		natsServer: natsServer,
	}

	for _, opt := range opts {
		opt(i)
	}

	nc, err := nats.Connect(i.natsServer)
	if err != nil {
		log.Fatal(err)
	}
	i.nc = nc

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	i.ec = ec

	return i
}

func (c *StrategyClient) Close() {
	c.ec.Close()
	c.nc.Close()
}

func (c *StrategyClient) SendSignalEvent(e ingenium.SignalEvent) error {
	subject := strings.ToLower(fmt.Sprintf("%s.%s", ingenium.SignalEventType, e.Symbol))
	return c.ec.Publish(subject, e)
}

func (c *StrategyClient) SubscribeToSymbol(symbol string, callback func(*ingenium.DataEvent)) error {
	_, err := c.ec.Subscribe(fmt.Sprintf("%s.%s", ingenium.DataEventType, strings.ToLower(symbol)), func(d *ingenium.DataEvent) {
		callback(d)
	})

	return err
}

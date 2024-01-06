package portfolio

import (
	"fmt"
	"log"
	"os"
	"strings"

	ingenium "github.com/markwinter/ingenium/pkg"
	"github.com/nats-io/nats.go"
)

type PortfolioClient struct {
	natsServer string
	nc         *nats.Conn
	ec         *nats.EncodedConn
}

type PortfolioOption func(*PortfolioClient)

// WithNatsServer sets the nats-server address. Default is http://0.0.0.0:4022 or read from env var NATS_SERVER if set
func WithNatsServer(server string) PortfolioOption {
	return func(i *PortfolioClient) {
		i.natsServer = server
	}
}

func NewPortfolioClient(portfolio ingenium.Portfolio, strategies []string, opts ...PortfolioOption) *PortfolioClient {
	var natsServer string

	natsServer = os.Getenv("NATS_SERVER")
	if natsServer == "" {
		natsServer = nats.DefaultURL
	}

	i := &PortfolioClient{
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

	for _, strategy := range strategies {
		subject := fmt.Sprintf("%s.%s", ingenium.SignalEventType, strings.ToLower(strategy))

		if _, err := i.ec.Subscribe(subject, func(d *ingenium.DataEvent) {
		}); err != nil {
			log.Printf("failed to subscribe to strategy: %v", err)
		}
	}

	return i
}

func (i *PortfolioClient) Close() {
	i.ec.Close()
	i.nc.Close()
}

func (i *PortfolioClient) SendOrder(order ingenium.OrderEvent) {}

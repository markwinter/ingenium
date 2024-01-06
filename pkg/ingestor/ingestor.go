package ingestor

import (
	"fmt"
	"log"
	"os"
	"strings"

	ingenium "github.com/markwinter/ingenium/pkg"
	"github.com/nats-io/nats.go"
)

type IngestorClient struct {
	natsServer string
	nc         *nats.Conn
	ec         *nats.EncodedConn
}

type IngestorOption func(*IngestorClient)

// WithNatsServer sets the nats-server address. Default is http://0.0.0.0:4022
func WithNatsServer(server string) IngestorOption {
	return func(i *IngestorClient) {
		i.natsServer = server
	}
}

func NewIngestor(opts ...IngestorOption) *IngestorClient {
	var natsServer string

	natsServer = os.Getenv("NATS_SERVER")
	if natsServer == "" {
		natsServer = nats.DefaultURL
	}

	i := &IngestorClient{
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

func (i *IngestorClient) Close() {
	i.ec.Close()
	i.nc.Close()
}

func (i *IngestorClient) SendDataEvent(d ingenium.DataEvent) error {
	subject := strings.ToLower(fmt.Sprintf("%s.%s", ingenium.DataEventType, d.Symbol))
	return i.ec.Publish(subject, d)
}

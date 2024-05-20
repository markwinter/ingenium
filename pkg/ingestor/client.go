package ingestor

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	ingenium "github.com/markwinter/ingenium/pkg"
	"github.com/nats-io/nats.go"
	"github.com/segmentio/ksuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const name = "ingenium"

var (
	tracer = otel.Tracer(name)
	meter  = otel.Meter(name)
)

type IngestorClient struct {
	natsServer string
	nc         *nats.Conn
	ec         *nats.EncodedConn

	eventMeter   metric.Int64Counter
	otelShutdown func(context.Context) error
}

type IngestorOption func(*IngestorClient)

// WithNatsServer sets the nats-server address. Default is http://0.0.0.0:4022
func WithNatsServer(server string) IngestorOption {
	return func(i *IngestorClient) {
		i.natsServer = server
	}
}

func NewIngestorClient(opts ...IngestorOption) *IngestorClient {
	ctx := context.Background()

	// Set up OpenTelemetry.
	otelShutdown, err := ingenium.SetupOTelSDK(ctx)
	if err != nil {
		return nil
	}

	var natsServer string

	natsServer = os.Getenv("NATS_SERVER")
	if natsServer == "" {
		natsServer = nats.DefaultURL
	}

	em, _ := meter.Int64Counter("event.data.sent")

	i := &IngestorClient{
		natsServer:   natsServer,
		eventMeter:   em,
		otelShutdown: otelShutdown,
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

	_ = i.otelShutdown(context.Background())
}

func (i *IngestorClient) SendDataEvent(e ingenium.DataEvent) error {
	ctx, span := tracer.Start(context.Background(), "event.data.send")
	defer span.End()

	eventAttr := attribute.String("event_id", e.Id)
	span.SetAttributes(eventAttr)

	i.eventMeter.Add(ctx, 1, metric.WithAttributes(eventAttr))

	e.Timestamp = time.Now()
	e.Id = "data_" + ksuid.New().String()

	subject := strings.ToLower(fmt.Sprintf("%s.%s", ingenium.DataEventType, e.Symbol))
	return i.ec.Publish(subject, e)
}

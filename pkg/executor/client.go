package executor

import (
	"log"
	"os"
	"time"

	ingenium "github.com/markwinter/ingenium/pkg"
	"github.com/nats-io/nats.go"
	"github.com/segmentio/ksuid"
)

type ExecutorClient struct {
	natsServer string
	nc         *nats.Conn
	ec         *nats.EncodedConn
}

type ExecutorOption func(*ExecutorClient)

// WithNatsServer sets the nats-server address. Default is http://0.0.0.0:4022
func WithNatsServer(server string) ExecutorOption {
	return func(i *ExecutorClient) {
		i.natsServer = server
	}
}

func NewExecutorClient(executor ingenium.Executor, opts ...ExecutorOption) *ExecutorClient {
	var natsServer string

	natsServer = os.Getenv("NATS_SERVER")
	if natsServer == "" {
		natsServer = nats.DefaultURL
	}

	i := &ExecutorClient{
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

	if _, err := i.ec.Subscribe(ingenium.OrderEventType, func(event *ingenium.OrderEvent) {
		executor.ReceiveOrder(event)
	}); err != nil {
		log.Printf("failed to subscribe to strategy signals: %v", err)
	}

	return i
}

func (i *ExecutorClient) Close() {
	i.ec.Close()
	i.nc.Close()
}

func (i *ExecutorClient) SendExecutionEvent(e ingenium.ExecutionEvent) error {
	e.Timestamp = time.Now()
	e.Id = "exec_" + ksuid.New().String()

	subject := ingenium.ExecutionEventType
	return i.ec.Publish(subject, e)
}

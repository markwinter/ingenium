package main

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

type EventPrinter struct {
	nc *nats.Conn
	ec *nats.EncodedConn
}

func MakeEventPrinter() EventPrinter {
	ep := EventPrinter{}

	natsServer := os.Getenv("NATS_SERVER")
	if natsServer == "" {
		natsServer = nats.DefaultURL
	}

	nc, err := nats.Connect(natsServer)
	if err != nil {
		log.Fatal(err)
	}
	ep.nc = nc

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	ep.ec = ec

	_, err = ec.Subscribe("ingenium.>", func(msg interface{}) {
		log.Printf("%v", msg)
	})

	if err != nil {
		log.Printf("printer failed to subscribe: %v", err)
	}

	return ep
}

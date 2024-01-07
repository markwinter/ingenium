package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	ingenium "github.com/markwinter/ingenium/pkg"
	"github.com/nats-io/nats.go"
)

var (
	subject string
)

func main() {
	flag.StringVar(&subject, "subject", "", "subject to subscribe and print")
	flag.Parse()

	if subject == "" {
		log.Fatal("required to set flag -subject")
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer ec.Close()

	_, err = ec.Subscribe(subject, func(data *ingenium.DataEvent) {
		log.Printf("%v\n", data)
	})

	if err != nil {
		log.Fatalf("failed to subscribe to subject: %v", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}

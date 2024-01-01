package main

import (
	"context"
	"fmt"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func receive(event cloudevents.Event) {
	fmt.Printf("%s", event)
}

func main() {
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Printf("Listening for events...\n")

	log.Fatal(c.StartReceiver(context.Background(), receive))
}

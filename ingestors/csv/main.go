package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	ingenium "github.com/markwinter/ingenium/pkg"
	"log"
	"os"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

func main() {
	csvPath := flag.String("csv", "", "Path to csv file")
	broker := flag.String("broker", "", "URL of broker to send events to")
	symbol := flag.String("symbol", "", "Security symbol")

	flag.Parse()

	if *csvPath == "" {
		log.Fatalf("csv file path not given")
	}

	if *broker == "" {
		log.Fatalf("broker url was not given")
	}

	if *symbol == "" {
		log.Fatalf("symbol was not given")
	}

	log.Printf("Reading csv file %s", *csvPath)
	records := readCsv(*csvPath)
	log.Printf("Read %d records", len(records))

	data := parseRecords(records, *symbol)

	log.Printf("Sending events to %s", *broker)
	sendEvents(data, *broker)

	log.Print("Finished")
}

func sendEvents(data []ingenium.DataEvent, broker string) {
	client, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	for i := 0; i < len(data); i++ {
		event := cloudevents.NewEvent()
		event.SetID(uuid.New().String())
		event.SetTime(time.Now())
		event.SetSource(fmt.Sprintf("ingenium/ingestor/csv/%s", os.Getenv("HOSTNAME")))
		event.SetType("ingenium.ingestor.data")

		event.SetData(cloudevents.ApplicationJSON, data[i])

		ctx := cloudevents.ContextWithTarget(context.Background(), broker)

		if result := client.Send(ctx, event); cloudevents.IsUndelivered(result) {
			log.Printf("failed to send, %v", result)
		}
	}
}

func parseRecords(records [][]string, symbol string) []ingenium.DataEvent {
	var data []ingenium.DataEvent
	for i := 1; i < len(records); i++ { // start at 1 to skip header row
		data = append(data, ingenium.DataEvent{
			Symbol:     symbol,
			Period:     records[i][0],
			OpenPrice:  records[i][1],
			ClosePrice: records[i][4],
			MaxPrice:   records[i][2],
			MinPrice:   records[i][3],
			Volume:     records[i][6],
		})
	}

	return data
}

func readCsv(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records
}

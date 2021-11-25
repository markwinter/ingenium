package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

type Data struct {
	Period     string
	OpenPrice  string
	ClosePrice string
	MaxPrice   string
	MinPrice   string
	Volume     string
}

func main() {
	csvPath := flag.String("csv", "", "Path to csv file")
	broker := flag.String("broker", "", "URL of broker to send events to")

	flag.Parse()

	if *csvPath == "" {
		log.Fatalf("csv file path not given")
	}

	if *broker == "" {
		log.Fatalf("broker url was not given")
	}

	log.Printf("Reading csv file %s", *csvPath)
	data := readCsv(*csvPath)
	log.Printf("Read %d records", len(data))

	log.Printf("Sending events to %s", *broker)
	sendEvents(data, *broker)

	log.Print("Finished")
}

func sendEvents(data []Data, broker string) {
	client, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	for i := 0; i < len(data); i++ {
		event := cloudevents.NewEvent()
		event.SetID(uuid.New().String())
		event.SetTime(time.Now())
		event.SetSource(fmt.Sprintf("ingenium/ingestor/csv/%s", os.Getenv("HOSTAME")))
		event.SetType("ingenium.ingestor.data")

		event.SetData(cloudevents.ApplicationJSON, map[string][]Data{"data": data})

		ctx := cloudevents.ContextWithTarget(context.Background(), broker)

		if result := client.Send(ctx, event); cloudevents.IsUndelivered(result) {
			log.Printf("failed to send, %v", result)
		}
	}
}

func readCsv(path string) []Data {
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

	var data []Data
	for i := 1; i < len(records); i++ {
		date, err := time.Parse("2006-01-02", records[i][0])
		if err != nil {
			log.Printf("Invalid date: %v\n%v", err, records[i])
			continue
		}

		data = append(data, Data{
			Period:     date.String(),
			OpenPrice:  records[i][1],
			ClosePrice: records[i][4],
			MaxPrice:   records[i][2],
			MinPrice:   records[i][3],
			Volume:     records[i][6],
		})
	}

	return data
}

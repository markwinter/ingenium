package ingenium

import (
	"context"
	"encoding/csv"
	"log"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// SendCsvFile is a helper function that reads a CSV file containing OHLC data and sends all the events to the broker.
// You must also supply a function that converts a row of csv data to an ingenium.DataEvent
func SendCsvFile(file string, sourceName, broker string, converter func(record []string) DataEvent) error {
	records, err := readCsv(file)
	if err != nil {
		return err
	}

	client, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	for _, record := range records {
		data := converter(record)
		event, err := ConvertDataEvent(data, sourceName)
		if err != nil {
			log.Printf("failed to convert data to cloudevent: %v", err)
		}

		ctx := cloudevents.ContextWithTarget(context.Background(), broker)
		if result := client.Send(ctx, event); cloudevents.IsUndelivered(result) {
			log.Printf("failed to send, %v", result)
		}
	}

	return nil
}

func readCsv(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records, nil
}

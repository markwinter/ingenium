package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	ingenium "github.com/markwinter/ingenium/pkg"
)

var (
	sourceName = fmt.Sprintf("ingenium/ingestor/csv/%s", os.Getenv("HOSTNAME"))
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

	err := ingenium.SendCsvFile(*csvPath, sourceName, *broker, func(record []string) ingenium.DataEvent {
		return ingenium.DataEvent{
			Type:      ingenium.DataTypeOhlc,
			Symbol:    *symbol,
			Timestamp: record[0],
			Data: ingenium.DataOhlc{
				Open:   record[1],
				High:   record[2],
				Low:    record[3],
				Close:  record[4],
				Volume: record[6],
				Period: "24h",
			},
		}
	})

	if err != nil {
		log.Fatalf("failed to send csv file: %v", err)
	}
}

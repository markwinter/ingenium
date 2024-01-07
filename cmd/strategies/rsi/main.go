package main

import (
	"flag"
	"log"

	"github.com/markwinter/ingenium/examples/strategies/rsi"
)

func main() {
	symbol := flag.String("symbol", "", "Security symbol")
	flag.Parse()

	if *symbol == "" {
		log.Fatalf("must set -symbol")
	}

	strategy := rsi.NewRsiStrategy(*symbol)
	defer strategy.Cleanup()

	strategy.Run()
}

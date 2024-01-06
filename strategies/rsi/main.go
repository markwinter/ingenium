package main

import (
	"flag"
	"log"

	rsi "github.com/markwinter/ingenium/strategies/rsi/pkg"
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

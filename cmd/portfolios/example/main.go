package main

import (
	"flag"

	exampleportfolio "github.com/markwinter/ingenium/examples/portfolios/example"
)

func main() {
	b := flag.Float64("balance", 10000.0, "Starting balance")

	flag.Parse()

	portfolio := exampleportfolio.NewPortfolio(*b)
	defer portfolio.Cleanup()

	portfolio.Run()
}

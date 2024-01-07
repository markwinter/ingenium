package alpacahistorical

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	ingenium "github.com/markwinter/ingenium/pkg"
	"github.com/markwinter/ingenium/pkg/ingestor"
)

const (
	feed = "IEX"
)

type AlpacaHistoricalIngestor struct {
	*ingestor.IngestorClient

	symbol    string
	timeframe marketdata.TimeFrame
	startDate time.Time
	endDate   time.Time
}

func NewAlpacaHistoricalIngestor(symbol string, startDate, endDate time.Time, timeframeString string) *AlpacaHistoricalIngestor {
	timeframe, err := convertTimeframe(timeframeString)
	if err != nil {
		log.Fatalf("failed to parse timeframe: %v", err)
	}

	i := &AlpacaHistoricalIngestor{
		IngestorClient: ingestor.NewIngestorClient(),
		symbol:         symbol,
		startDate:      startDate,
		endDate:        endDate,
		timeframe:      timeframe,
	}

	return i
}

func (i AlpacaHistoricalIngestor) IngestData() {
	client := marketdata.NewClient(marketdata.ClientOpts{
		Feed: feed,
	})

	bars, err := client.GetBars(i.symbol, marketdata.GetBarsRequest{
		TimeFrame: i.timeframe,
		Start:     i.startDate,
		End:       i.endDate,
	})
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%s bars:\n", i.symbol)
	for _, bar := range bars {
		d := i.convertToDataEvent(bar)
		//fmt.Printf("%+v\n", d)

		if err := i.SendDataEvent(d); err != nil {
			log.Printf("failed sending data event: %v", err)
		}
	}
}

func (i AlpacaHistoricalIngestor) Cleanup() {
	i.Close()
}

func (i *AlpacaHistoricalIngestor) convertToDataEvent(bar marketdata.Bar) ingenium.DataEvent {
	return ingenium.DataEvent{
		Type:      ingenium.DataTypeOhlc,
		Timestamp: time.Now(),
		Symbol:    i.symbol,
		Ohlc: ingenium.DataOhlc{
			Open:      fmt.Sprintf("%f", bar.Open),
			High:      fmt.Sprintf("%f", bar.High),
			Low:       fmt.Sprintf("%f", bar.Low),
			Close:     fmt.Sprintf("%f", bar.Close),
			Volume:    fmt.Sprintf("%d", bar.Volume),
			Period:    i.timeframe.String(),
			Timestamp: bar.Timestamp,
		},
	}
}

// convertTimeframe converts a string like "30m" to an alpaca marketdata.TimeFrame
func convertTimeframe(timeframe string) (marketdata.TimeFrame, error) {
	var n strings.Builder
	var u rune

	// Parse a string like "30m" into n=30, u=m
	for _, c := range timeframe {
		if !unicode.IsNumber(c) {
			u = c
			break
		}

		n.WriteRune(c)
	}

	time, err := strconv.Atoi(n.String())
	if err != nil {
		return marketdata.TimeFrame{}, err
	}

	var unit marketdata.TimeFrameUnit
	switch u {
	case 'h':
		unit = marketdata.Hour
	case 'm':
		unit = marketdata.Min
	case 'd':
		unit = marketdata.Day
	default:
		return marketdata.TimeFrame{}, fmt.Errorf("unknown time unit %q", u)
	}

	return marketdata.NewTimeFrame(time, unit), nil
}

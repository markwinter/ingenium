package ingenium

type Ingestor interface {
	IngestData()
	Cleanup()
}

type Strategy interface {
	Run()
	Receive(*DataEvent)
	Cleanup()
}

type Executor interface {
	Run()
	Cleanup()
}

type Portfolio interface {
	Run()
	Cleanup()
}

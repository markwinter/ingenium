package ingenium

type Ingestor interface {
	Cleanup()
	IngestData()
}

type Strategy interface {
	Run()
	Cleanup()
	Receive(*DataEvent)
}

type Executor interface {
	Run()
	Cleanup()
	ReceiveOrder(*OrderEvent)
}

type Portfolio interface {
	Run()
	Cleanup()
	ReceiveSignal(*SignalEvent)
	ReceiveExecution(*ExecutionEvent)
}

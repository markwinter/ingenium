package ingenium

import (
	"fmt"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/segmentio/ksuid"
)

// ConvertDataEvent converts an Ingenium DataEvent to a CloudEvent
func ConvertDataEvent(data DataEvent, sourceName string) (cloudevents.Event, error) {
	id := fmt.Sprintf("data_%s", ksuid.New())

	event := cloudevents.NewEvent()
	event.SetID(id)
	event.SetTime(time.Now())
	event.SetSource(sourceName)
	event.SetType(DataEventType)

	if err := event.SetData(cloudevents.ApplicationJSON, data); err != nil {
		return cloudevents.Event{}, err
	}

	return event, nil
}

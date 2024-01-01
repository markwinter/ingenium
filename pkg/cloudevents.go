package ingenium

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

// ConvertDataEvent converts an Ingenium DataEvent to a CloudEvent
func ConvertDataEvent(data DataEvent, sourceName string) (cloudevents.Event, error) {
	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetTime(time.Now())
	event.SetSource(sourceName)
	event.SetType(DataEventType)

	if err := event.SetData(cloudevents.ApplicationJSON, data); err != nil {
		return cloudevents.Event{}, err
	}

	return event, nil
}

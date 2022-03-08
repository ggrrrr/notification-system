package eventbus

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Will try and publish the event to event.Topic with partitional Key
func (ep *eventProducer) PublishEvent(event *Event) {
	var payload []byte
	var err error
	if ep.producer == nil {
		logrus.Errorf("producer is nil")
		return
	}
	if event.Topic == "" {
		logrus.Errorf("event.Topic is nil")
		return
	}
	if event.RecordKey == "" {
		event.RecordKey = uuid.New().String()
	}
	if event.message == nil {
		if event.Payload == nil {
			logrus.Errorf("missing event.payload")
			return
		}
		event.ContentType = H_CT_JSON
		payload, err = json.Marshal(event.Payload)
		if err != nil {
			logrus.Error(err)
			return
		}
		event.message = payload
	}
	ep.push(event)
}

func (ep *eventProducer) PublishRetry(event *Event) {
	event.RetryCount--
	// Making sure we dont loop
	if event.RetryCount == 0 {
		event.Error = "RetryCount reached 0"
		ep.PublishError(event)
		return
	}
	event.SenderTopic = event.Topic
	event.Topic = retryTopic
	event.RecordKey = uuid.New().String()
	ep.push(event)
}

func (ep *eventProducer) PublishDone(event *Event) {
	event.Topic = doneTopic
	event.SenderTopic = event.Topic
	ep.push(event)
}

func (ep *eventProducer) PublishError(event *Event) {
	event.Topic = errorTopic
	event.SenderTopic = event.Topic
	ep.push(event)
}

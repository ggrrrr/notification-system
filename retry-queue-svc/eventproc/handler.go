package eventproc

import (
	"time"

	"github.com/ggrrrr/notification-system/common-lib/eventbus"
	"github.com/sirupsen/logrus"
)

// Process single event
func EventHandler(event *eventbus.Event) {
	logrus.Debugf("event: app: %v counter: %v, sender: %v, key: %v",
		event.Sender, event.RetryCount, event.SenderTopic, event.RecordKey)
	if event.SenderTopic == "" {
		logrus.WithFields(logrus.Fields{
			"event": event,
		}).Errorf("SenderTopic is missing")
		return
	}
	if event.RetryCount < 1 {
		logrus.WithFields(logrus.Fields{
			"event": event,
		}).Errorf("RetryCount < 1")
		return
	}
	if event.RetryCount > eventbus.GetRetryCounter() {
		logrus.WithFields(logrus.Fields{
			"event": event,
		}).Errorf("RetryCount > %v", eventbus.GetRetryCounter())
		return
	}

	// event.RetryCount--
	// logrus.WithField("event", event).Debugf("queue: %v", msg.TopicPartition)
	eventC.Pause(event)
	time.Sleep(sleepTime)
	// logrus.WithField("event", event).Debugf("resume: %v", msg.TopicPartition)
	eventC.Resume(event)
	event.Topic = event.SenderTopic
	eventP.PublishEvent(event)
}

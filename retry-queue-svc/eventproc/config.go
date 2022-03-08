package eventproc

import (
	"fmt"
	"time"

	"github.com/ggrrrr/notification-system/common-lib/app"
	"github.com/ggrrrr/notification-system/common-lib/config"
	"github.com/ggrrrr/notification-system/common-lib/eventbus"
	"github.com/sirupsen/logrus"
)

const (
	C_PREFIX = "retry.queue"
	SLEEP    = "sleep.time"
)

var (
	sleepTime time.Duration
	topic     string
	err       error
	eventC    eventbus.EventConsumer
	eventP    eventbus.EventProducer
	eventChan = make(chan *eventbus.Event)
)

// Configure the consumer of notification_reuqest
func Configure() error {
	eventP, err = eventbus.NewProducer()
	if err != nil {
		return err
	}

	topic = eventbus.GetRetryTopic()
	sleepTime = config.GetDuration(C_PREFIX, SLEEP, 30*time.Second)

	if topic == "" {
		logrus.Infof("topic: %v", topic)
		return fmt.Errorf("retry.queue.topic not set")
	}

	eventbus.CreateTopic(topic, 1)

	eventC, err = eventbus.NewConsumer("eventbus.retry")
	if err != nil {
		return err
	}
	return nil
}

// starts consumer
func Start() {
	eventC.Start(eventChan)
}

// Blocking code Starts the proccessing loop
func EventLoop() {
	for !app.ShuttingDown() {
		event := <-eventChan
		EventHandler(event)
	}
	logrus.Infof("end.")
}

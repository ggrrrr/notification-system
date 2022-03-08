package eventproc

import (
	"fmt"

	"github.com/ggrrrr/notification-system/common-lib/app"
	"github.com/ggrrrr/notification-system/common-lib/config"
	"github.com/ggrrrr/notification-system/common-lib/eventbus"
	"github.com/sirupsen/logrus"
)

const (
	C_PREFIX = "event"
)

var (
	topic     string
	err       error
	eventC    eventbus.EventConsumer
	eventP    eventbus.EventProducer
	eventChan = make(chan *eventbus.Event)
)

// Configure the consumer of notification_reuqest
func Configure() error {
	topic = config.GetString(C_PREFIX, "topic")
	if topic == "" {
		return fmt.Errorf("invalid topic")
	}
	eventP, err = eventbus.NewProducer()
	if err != nil {
		return err
	}
	eventbus.CreateTopic(topic, 1)

	eventC, err = eventbus.NewConsumer(C_PREFIX)
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

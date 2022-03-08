package notification

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type register struct {
	channels map[string]NotificationService
}

var r *register

// init the channel register
func Init() {
	r = &register{
		channels: map[string]NotificationService{},
	}
}

// Add channel to the registery
func Add(channel NotificationService) {
	logrus.Infof("add: %v", channel.Name())
	r.channels[channel.Name()] = channel
}

// Attempt to push notification request to a channel
func Process(msg *NotificationData) (EVENT_RESULT, error) {
	logrus.Infof("msg: channel:%v", msg.Channel)
	c, ok := r.channels[msg.Channel]
	if !ok {
		return EVENT_ERROR, fmt.Errorf("unkown channel: %v", msg.Channel)
	}
	forRetry, err := c.Push(msg)

	return forRetry, err
}

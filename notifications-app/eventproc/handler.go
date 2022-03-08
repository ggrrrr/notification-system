package eventproc

import (
	"github.com/ggrrrr/notification-system/common-lib/eventbus"
	"github.com/ggrrrr/notification-system/notifications-app/notification"
	"github.com/sirupsen/logrus"
)

// Process single event
func EventHandler(event *eventbus.Event) {
	logrus.Debugf("event: %v", event)
	var payload notification.NotificationData
	err := event.Unmarshal(&payload)
	if err != nil {
		logrus.Errorf("Unmarshal: %v", err)
		return
	}
	retry, err := notification.Process(&payload)
	if err != nil {
		event.Error = err.Error()
	}
	switch retry {
	case notification.EVENT_RETRY:
		logrus.Infof("event retry: %v", err)
		eventP.PublishRetry(event)
	case notification.EVENT_ERROR:
		logrus.Infof("event error: %v", err)
		eventP.PublishError(event)
	case notification.EVENT_DONE:
		logrus.Infof("event done")
		eventP.PublishDone(event)
	}

}

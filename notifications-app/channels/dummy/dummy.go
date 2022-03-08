package dummy

import (
	"fmt"

	"github.com/ggrrrr/notification-system/common-lib/config"
	"github.com/ggrrrr/notification-system/notifications-app/notification"
	"github.com/sirupsen/logrus"
)

type cfg struct {
	name string
}

// Creates and configures Sms channel
func New() (notification.NotificationService, error) {
	c := cfg{name: "dummy"}
	err := c.configure()
	return &c, err
}

func (c *cfg) Push(msg *notification.NotificationData) (notification.EVENT_RESULT, error) {
	logrus.Infof("msg: %v", msg)
	d, ok := msg.Body["dummy"]
	if !ok {
		return notification.EVENT_ERROR, fmt.Errorf("bad body for channel: %v", c.name)
	}
	if d == "retry" {
		return notification.EVENT_RETRY, fmt.Errorf("dummy:%v", d)
	}
	if d == "error" {
		return notification.EVENT_ERROR, fmt.Errorf("dummy:%v", d)
	}
	return notification.EVENT_DONE, nil
}

func (c *cfg) Name() string {
	return c.name
}

func (c *cfg) configure() error {
	if !config.GetBool(c.name, notification.C_ENABLE) {
		return fmt.Errorf("not enabled")
	}

	return nil
}

package slack

import (
	"fmt"

	"github.com/ggrrrr/notification-system/common-lib/config"
	"github.com/ggrrrr/notification-system/notifications-app/notification"
	"github.com/sirupsen/logrus"
)

type cfg struct {
	name  string
	url   string
	token string
}

// Creates and configures Slack channel
func New() (notification.NotificationService, error) {
	c := cfg{name: "slack"}
	err := c.configure()
	return &c, err
}

func (c *cfg) Name() string {
	return c.name
}

func (c *cfg) configure() error {
	if !config.GetBool(c.name, "enable") {
		return fmt.Errorf("not enabled")
	}
	c.token = config.GetString(c.name, "api.token")
	c.url = config.GetString(c.name, "api.url")
	if c.url == "" {
		return fmt.Errorf("api.url not set")
	}
	if c.token == "" {
		return fmt.Errorf("api.token not set")
	}
	logrus.Infof("url: %v", c.url)
	return nil
}

package email

import (
	"fmt"

	"github.com/ggrrrr/notification-system/common-lib/config"
	"github.com/ggrrrr/notification-system/notifications-app/notification"
	"github.com/sirupsen/logrus"
)

type cfg struct {
	name       string
	smtpServer string
	smtpPort   int
	username   string
	password   string
	sender     string
	tls        bool
}

// Creates and configures Sms channel
func New() (notification.NotificationService, error) {
	c := cfg{name: "email"}
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
	c.smtpServer = config.GetString(c.name, "smtp.server")
	if c.smtpServer == "" {
		return fmt.Errorf("smtp.server not set")
	}
	c.smtpPort = config.GetInt(c.name, "smtp.port", 0)
	if c.smtpPort == 0 {
		return fmt.Errorf("smtp.port not set")
	}
	c.username = config.GetString(c.name, "smtp.username")
	if c.username == "" {
		return fmt.Errorf("smtp.username not set")
	}
	c.password = config.GetString(c.name, "smtp.password")
	if c.password == "" {
		return fmt.Errorf("smtp.password not set")
	}
	c.sender = config.GetString(c.name, "smtp.sender")
	if c.password == "" {
		return fmt.Errorf("smtp.sender not set")
	}
	c.tls = config.GetBool(c.name, "smtp.tls")
	logrus.Infof("smtpServer: %v:%v", c.smtpServer, c.smtpPort)
	return nil
}

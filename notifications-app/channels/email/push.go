package email

import (
	"fmt"
	"net/smtp"

	"github.com/ggrrrr/notification-system/common-lib/config"
	"github.com/ggrrrr/notification-system/notifications-app/notification"
)

func (c *cfg) Push(msg *notification.NotificationData) (notification.EVENT_RESULT, error) {
	auth := smtp.PlainAuth("", c.username, c.password, c.smtpServer)
	sender := c.sender
	smtpHostPort := fmt.Sprintf("%v:%v", c.smtpServer, c.smtpPort)

	to1, ok := msg.Body["to"]
	if !ok {
		return notification.EVENT_ERROR, fmt.Errorf("missing to")
	}
	subject1, ok := msg.Body["subject"]
	if !ok {
		return notification.EVENT_ERROR, fmt.Errorf("missing subject")
	}
	text1, ok := msg.Body["text"]
	if !ok {
		return notification.EVENT_ERROR, fmt.Errorf("missing text")
	}
	sender1, ok := msg.Body["sender"]
	if ok {
		sender = config.ItoS(sender1)
	}
	to := []string{config.ItoS(to1)}
	subject := config.ItoS(subject1)
	text := config.ItoS(text1)
	templ := "To: %v\r\n" +
		"Subject: %v\r\n" +
		"\r\n" +
		"%v\r\n"
	body := fmt.Sprintf(templ, to, subject, text)
	err := smtp.SendMail(smtpHostPort, auth, sender, to, []byte(body))
	if err != nil {
		return notification.EVENT_RETRY, err
	}
	return notification.EVENT_DONE, nil
}

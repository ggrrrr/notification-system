package email_test

import (
	"fmt"
	"testing"

	smtpmock "github.com/mocktools/go-smtp-mock"

	"github.com/ggrrrr/notification-system/notifications-app/channels/email"
	"github.com/ggrrrr/notification-system/notifications-app/notification"
)

func TestMail(t *testing.T) {
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       true,
		LogServerActivity: true,
	})
	if err := server.Start(); err != nil {
		fmt.Println(err)
	}
	defer server.Stop()
	port := fmt.Sprintf("%v", server.PortNumber)
	t.Setenv("EMAIL_ENABLE", "true")
	t.Setenv("EMAIL_SMTP_SERVER", "localhost")
	t.Setenv("EMAIL_SMTP_PORT", port)
	t.Setenv("EMAIL_SMTP_USERNAME", "user")
	t.Setenv("EMAIL_SMTP_PASSWORD", "pass")
	t.Setenv("EMAIL_SMTP_SENDER", "sender")
	channel, err := email.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	msg := notification.NotificationData{
		Channel: "email",
		Body: map[string]interface{}{
			"text":    "text me",
			"to":      "asdasd",
			"subject": "new subject",
		},
	}
	ok, err := channel.Push(&msg)
	if ok != notification.EVENT_RETRY {
		t.Errorf("ok %v %v", ok, err)
	}

}

func TestMail12(t *testing.T) {

	t.Setenv("EMAIL_ENABLE", "true")
	t.Setenv("EMAIL_SMTP_SERVER", "localhost")
	t.Setenv("EMAIL_SMTP_PORT", "122")
	t.Setenv("EMAIL_SMTP_USERNAME", "user")
	t.Setenv("EMAIL_SMTP_PASSWORD", "pass")
	t.Setenv("EMAIL_SMTP_SENDER", "sender")
	channel, err := email.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	msg := notification.NotificationData{
		Channel: "email",
		Body:    map[string]interface{}{"text": "text me", "to": "asdasd", "subject": "new subject"},
	}
	ok, err := channel.Push(&msg)
	if ok != notification.EVENT_RETRY {
		t.Errorf("ok %v %v", ok, err)
	}

}

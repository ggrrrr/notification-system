package slack_test

import (
	"testing"

	"github.com/ggrrrr/notification-system/notifications-app/channels/slack"
	"github.com/ggrrrr/notification-system/notifications-app/notification"
)

func TestPush(t *testing.T) {
	slack, err := slack.New()
	if err != nil {
		t.Fatalf("%v", err)
	}

	msg := notification.NotificationData{
		Channel: "slack",
		Body: map[string]interface{}{
			"text":    "text me",
			"channel": "asdasd",
		},
	}
	ok, err := slack.Push(&msg)
	t.Logf("ok %v %v", ok, err)

}

package notification_test

import (
	"testing"

	"github.com/ggrrrr/notification-system/notifications-app/channels/dummy"
	"github.com/ggrrrr/notification-system/notifications-app/notification"
)

func TestRegistry(t *testing.T) {
	var err error
	var retry notification.EVENT_RESULT
	notification.Init()

	dummy, err := dummy.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	notification.Add(dummy)

	msg := notification.NotificationData{Channel: "dummy", Body: map[string]interface{}{"dummy": "ok"}}
	msg1 := notification.NotificationData{Channel: "dummy", Body: map[string]interface{}{"dumm1y": "some"}}
	msg2 := notification.NotificationData{Channel: "dummya", Body: map[string]interface{}{"dumm1y": "some"}}
	msg3 := notification.NotificationData{Channel: "dummy", Body: map[string]interface{}{"dummy": "retry"}}

	retry, err = notification.Process(&msg)
	if retry != notification.EVENT_DONE {
		t.Errorf("must be done: %v %v", retry, err)
	}

	retry, err = notification.Process(&msg1)
	if retry != notification.EVENT_ERROR {
		t.Errorf("must be error: %v %v", retry, err)
	}

	retry, err = notification.Process(&msg2)
	if retry != notification.EVENT_ERROR {
		t.Errorf("must be error: %v %v", retry, err)
	}
	retry, err = notification.Process(&msg3)
	if retry != notification.EVENT_RETRY {
		t.Errorf("must be error: %v %v", retry, err)
	}

}

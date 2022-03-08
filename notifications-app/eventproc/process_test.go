package eventproc_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ggrrrr/notification-system/common-lib/eventbus"
	"github.com/ggrrrr/notification-system/notifications-app/channels/dummy"
	"github.com/ggrrrr/notification-system/notifications-app/eventproc"
	"github.com/ggrrrr/notification-system/notifications-app/notification"
	"github.com/google/uuid"
)

var (
	topic    string
	err      error
	producer eventbus.EventProducer
)

func push(msg *notification.NotificationData) {
	e := eventbus.Event{
		Topic:  topic,
		Sender: "test-app",

		Payload: msg,
	}
	producer.PublishEvent(&e)
}

func TestEventBus(t *testing.T) {
	ttt := uuid.New().String()[0:4]
	topic = "notification_request" + ttt
	t.Setenv("DUMMY_ENABLE", "true")
	t.Setenv("EVENT_TOPIC", topic)
	t.Setenv("EVENT_GROUP_ID", "mygroup")
	t.Setenv("EVENT_AUTO_OFFSET_RESET", "beginning")

	var wg sync.WaitGroup
	wg.Add(1)

	producer, err = eventbus.NewProducer()
	if err != nil {
		t.Fatal(err)
	}

	msg1 := notification.NotificationData{
		Channel: "dummy",
		Body:    map[string]interface{}{"dummy": "text me"},
	}
	msg2 := notification.NotificationData{
		Channel: "dummy",
		Body:    map[string]interface{}{"dummy": "retry"},
	}

	notification.Init()
	d, err := dummy.New()
	if err != nil {
		t.Error(err)
	}
	notification.Add(d)
	// ok, err := notification.Process(&msg1)
	// t.Logf(" %v %v", ok, err)
	// ok, err = notification.Process(&msg2)
	// t.Logf(" %v %v", ok, err)

	err = eventproc.Configure()
	if err != nil {
		panic(err)
	}
	push(&msg1)
	push(&msg2)
	time.Sleep(10 * time.Second)

	eventproc.Start()

	wg.Wait()
}

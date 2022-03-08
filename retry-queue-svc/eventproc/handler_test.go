package eventproc_test

import (
	"testing"
	"time"

	"github.com/ggrrrr/notification-system/common-lib/eventbus"
	"github.com/ggrrrr/notification-system/retry-queue-svc/eventproc"
	"github.com/google/uuid"
)

type TestData struct {
	Str string `json:"str"`
}

var err error

func Test1(t *testing.T) {

	err = eventproc.Configure()
	if err != nil {
		t.Errorf("NewProducer config %v", err)
	}
	eventproc.Start()

	data := TestData{Str: uuid.New().String()}
	event1 := eventbus.Event{
		Topic:       "test",
		SenderTopic: "notification_request",
		ContentType: eventbus.H_CT_JSON,
		RetryCount:  3,
		Payload:     data,
	}

	eventproc.EventHandler(&event1)
	time.Sleep(10000)

}

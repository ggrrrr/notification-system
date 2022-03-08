package eventbus_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ggrrrr/notification-system/common-lib/eventbus"
	"github.com/google/uuid"
)

var wg sync.WaitGroup

func Test1(t *testing.T) {
	var err error
	ttt := uuid.New().String()[0:4]

	topic := "test_topic_test" + ttt
	t.Setenv("TEST_TOPIC", topic)
	// t.Setenv("TEST_AUTO_OFFSET_RESET", "smallest")

	producer, err := eventbus.NewProducer()
	if err != nil {
		t.Fatalf("producer %v", err)
	}
	eventbus.CreateTopic(topic, 1)
	time.Sleep(10 * time.Second)
	payload := struct{ Key string }{Key: uuid.New().String()}
	testEvent1 := eventbus.NewEvent(topic, payload)
	// var testEvent2 eventbus.Event
	events := make(chan *eventbus.Event)
	testC, err := eventbus.NewConsumer("test")
	if err != nil {
		t.Fatalf("test consumer config %v", err)
	}
	testC.Start(events)
	// testC.WaitAssign()
	producer.PublishEvent(testEvent1)
	if err != nil {
		t.Errorf("producer %v", err)
	}
	wg.Add(1)
	wg.Add(1)
	recived := true
	for recived {
		testEvent3 := <-events
		// t.Logf("testEvent1 %+v", testEvent1)
		// t.Logf("testEvent3 %+v", testEvent3)
		t.Logf("testEvent1 %+v", string(testEvent1.EventHash))
		t.Logf("testEvent3 %+v", string(testEvent3.EventHash))
		wg.Done()
		if testEvent3.EventHash == testEvent1.EventHash {
			recived = false
		}

	}
	producer.WaitAll()
	testC.Close()
}

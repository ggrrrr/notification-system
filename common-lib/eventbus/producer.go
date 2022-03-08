package eventbus

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type eventProducer struct {
	producer *kafka.Producer
	wg       sync.WaitGroup
}

var (
	producer *kafka.Producer
)

func NewProducer() (EventProducer, error) {
	if !hasString(kafkaConfig["bootstrap.servers"]) {
		return nil, fmt.Errorf("bootstrap.servers not set")
	}

	p, err := kafka.NewProducer(&kafkaConfig)
	if err != nil {
		logrus.WithField("config", kafkaConfig).Errorf("unable to create producer: %v", err)
		return nil, err
	}
	out := eventProducer{
		producer: p,
		wg:       sync.WaitGroup{},
	}
	go out.producerEventsHandler()
	producer = p
	return &out, nil
}

func (ep *eventProducer) Close() {
	ep.producer.Close()
}

func (ep *eventProducer) producerEventsHandler() {
	for e := range ep.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			ep.wg.Done()
			if ev.TopicPartition.Error != nil {
				logrus.WithFields(logrus.Fields{
					"topic":     *ev.TopicPartition.Topic,
					"partition": ev.TopicPartition.Partition,
				}).Errorf("Failed: message: %v", string(ev.Value))
			} else {
				logrus.WithFields(logrus.Fields{
					"topic":     *ev.TopicPartition.Topic,
					"partition": ev.TopicPartition.Partition,
				}).Debugf("published: message: %v", string(ev.Value))
			}
		default:
			logrus.Errorf("unkown")
		}
	}
}

func (ep *eventProducer) push(event *Event) {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &event.Topic, Partition: kafka.PartitionAny},
		Value:          event.message,
		Opaque:         nil,
	}

	message.Headers = createHeaders(event, event.message)
	message.Key = []byte(event.RecordKey)
	logrus.WithFields(logrus.Fields{
		"senderTopic": event.SenderTopic,
		"recordKey":   event.RecordKey,
		"retryCount":  event.RetryCount,
		"hash":        event.EventHash,
	}).Debugf("metadata")
	logrus.Debugf("payload: %v: %v", event.Topic, string(event.message))
	err := ep.producer.Produce(message, nil)
	if err != nil {
		logrus.Error(err)
		return
	}
	ep.wg.Add(1)
}

// Wait all pending to be clear or timeout
func (ep *eventProducer) WaitAll() {
	ep.wg.Wait()
}

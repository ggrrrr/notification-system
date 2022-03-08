package eventbus

import (
	"fmt"

	"github.com/ggrrrr/notification-system/common-lib/app"
	"github.com/ggrrrr/notification-system/common-lib/config"
	"github.com/sirupsen/logrus"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type eventConsumer struct {
	topic            string
	consumer         *kafka.Consumer
	eventChan        chan *Event
	msgCounter       int
	minCommitCounter int
}

func NewConsumer(prefix string) (EventConsumer, error) {
	// func NewConsumer(prefix string, sd chan struct{}) (*eventConsumer, error) {
	cfg := consumerConfig(prefix)

	topic := config.GetString(prefix, "topic")
	commmitCouner := config.GetInt(prefix, "commit.counter", 5)
	allCfg := copyConfig(kafkaConfig, cfg)

	if !hasString(allCfg["bootstrap.servers"]) {
		return nil, fmt.Errorf("bootstrap.servers not set")
	}

	if topic == "" {
		return nil, fmt.Errorf("NewConsumer topic not set")
	}

	// allCfg["go.events.channel.enable"] = true
	// allCfg["go.application.rebalance.enable"] = true
	// allCfg["enable.auto.commit"] = true
	c, err := kafka.NewConsumer(&allCfg)
	if err != nil {
		logrus.WithField("config", allCfg).Errorf("unable to create consumer: %v", err)
		return nil, err
	}
	// TODO set config for min/max counter
	out := eventConsumer{
		topic:            topic,
		consumer:         c,
		msgCounter:       0,
		minCommitCounter: commmitCouner,
	}

	return &out, nil
}

func (ev *eventConsumer) Close() {
	ev.consumer.Close()
}

// Non blocking code
func (ec *eventConsumer) Start(eventChan chan *Event) error {
	err := ec.consumer.SubscribeTopics([]string{ec.topic}, nil)
	// err = c.SubscribeTopics([]string{topic}, out.rebalanceCb)
	if err != nil {
		logrus.Errorf("SubscribeTopics(%v): %v", ec.topic, err)
		return err
	}
	logrus.Infof("SubscribeTopics: %v", ec.topic)
	ec.eventChan = eventChan
	go ec.loop()
	return nil
}

// Blocking code: please use goroutines
func (ec *eventConsumer) loop() {
	logrus.Infof("start...")
	defer ec.consumer.Close()
	for !app.ShuttingDown() {
		ev := ec.consumer.Poll(1)
		switch e := ev.(type) {
		case kafka.AssignedPartitions:
			logrus.Infof("%v", e)
			// c.Assign(e.Partitions)
		case kafka.RevokedPartitions:
			logrus.Infof("%v", e)
			// c.Unassign()
		case *kafka.Message:
			logrus.Infof("%% Message on %s:%s", e.TopicPartition, string(e.Value))
			if ec.minCommitCounter > 0 {
				ec.msgCounter++
				if ec.msgCounter%ec.minCommitCounter == 0 {
					ec.consumer.Commit()
				}
			}
			event := parseMessage(e)
			if event.ContentType == "" {
				logrus.Warnf("event without content-type")
				continue
			}
			// TODO JWT or other Auth
			ec.eventChan <- event
		case kafka.PartitionEOF:
			logrus.Errorf("PartitionEOF Reached %v\n", e)
		case kafka.Error:
			logrus.Errorf("Error: %v", e)
			app.Panic(fmt.Errorf("kafka error: %v", e))
		default:
		}
	}
	ec.consumer.Commit()
	logrus.Info("shutdown")
}

func (ec *eventConsumer) rebalanceCb(consumer *kafka.Consumer, ev kafka.Event) error {
	switch e := ev.(type) {
	case kafka.AssignedPartitions:
		for _, v := range e.Partitions {
			logrus.Infof("Assigned: %v/%v", v.Topic, v.Partition)
		}
	case kafka.RevokedPartitions:
		for _, v := range e.Partitions {
			logrus.Infof("Revoked: %v/%v", v.Topic, v.Partition)
		}
	}
	return nil
}

func (ec *eventConsumer) Pause(event *Event) {
	err := ec.consumer.Pause([]kafka.TopicPartition{event.topicPartition})
	if err != nil {
		logrus.Error(err)
	}
}

func (ec *eventConsumer) Resume(event *Event) {
	err := ec.consumer.Resume([]kafka.TopicPartition{event.topicPartition})
	if err != nil {
		logrus.Error(err)
	}
}

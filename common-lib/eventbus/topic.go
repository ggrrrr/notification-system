package eventbus

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

// Will attempt to create topic
func CreateTopic(topic string, partitions int) {
	if producer == nil {
		logrus.Errorf("producer is null")
	}
	a, err := kafka.NewAdminClientFromProducer(producer)
	if err != nil {
		logrus.Errorf("Failed to create admin: %s", err)
	}
	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		logrus.Errorf("ParseDuration(60s): %s", err)
	}
	results, err := a.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{{
			Topic:         topic,
			NumPartitions: partitions,
		}},
		// Admin options
		kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		logrus.Errorf("Admin Client request error: %v", err)
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			logrus.Errorf("Failed to create topic: %v", result.Error)
		}
		logrus.Info(result)
	}
	a.Close()
}

package eventbus_test

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
)

func Test2(t *testing.T) {

	servers := os.Getenv("BOOTSTRAP_SERVERS")
	topic := "test_topic"

	var wg sync.WaitGroup
	wg.Add(2)

	random := uuid.New().String()
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
		"client.id":         "local1qwe",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					t.Logf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					wg.Done()
					t.Logf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	// delivery_chan := make(chan kafka.Event, 10000)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(random)},
		nil,
	)
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	t.Logf("%v", random)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
		"group.id":          "foo",
		"auto.offset.reset": "smallest"})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	err = consumer.SubscribeTopics([]string{topic}, nil)
	run := true
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%% Message on :\n%s", random)

	go func() {
		for run == true {
			ev := consumer.Poll(0)
			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				gotMsg := string(e.Value)
				if gotMsg == random {
					wg.Done()
				}
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			default:
				// fmt.Printf("Ignored %v\n", e)
			}
		}
	}()
	wg.Wait()
}

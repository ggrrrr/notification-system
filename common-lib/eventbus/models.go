package eventbus

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ggrrrr/notification-system/common-lib/app"
)

type Event struct {
	Sender         string
	Topic          string
	SenderTopic    string
	ContentType    string
	RetryCount     int
	EventHash      string
	RecordKey      string
	Payload        interface{}
	Error          string
	message        []byte
	topicPartition kafka.TopicPartition
	// headers     []kafka.Header
	// message     *kafka.Message
}

// Create new event for topic and event payload
func NewEvent(topic string, payload interface{}) *Event {

	out := Event{
		ContentType: H_CT_JSON,
		Sender:      app.GetName(),
		Topic:       topic,
		Payload:     payload,
		RetryCount:  retryCounter,
	}
	return &out
}

// Create new event from JSON for topic
func NewFromJson(topic string, message []byte, payload interface{}) (*Event, error) {
	err := json.Unmarshal(message, &payload)
	if err != nil {
		return nil, err
	}

	out := Event{
		ContentType: H_CT_JSON,
		Sender:      app.GetName(),
		Topic:       topic,
		Payload:     payload,
		message:     message,
		RetryCount:  retryCounter,
	}
	return &out, nil
}

func createHeaders(event *Event, payload []byte) []kafka.Header {
	hash := sha256.Sum256(payload)
	event.EventHash = base64.URLEncoding.EncodeToString(hash[:])
	out := []kafka.Header{}
	out = append(out, kafka.Header{Key: H_SENDER_APP, Value: []byte(event.Sender)})
	if event.ContentType != "" {
		out = append(out, kafka.Header{Key: H_CONTENT_TYPE, Value: []byte(event.ContentType)})
	}
	if event.Topic != "" {
		out = append(out, kafka.Header{Key: H_TOPIC, Value: []byte(event.Topic)})
	}
	if event.SenderTopic != "" {
		out = append(out, kafka.Header{Key: H_SENDER_TOPIC, Value: []byte(event.SenderTopic)})
	}
	if event.RecordKey != "" {
		out = append(out, kafka.Header{Key: H_RECORD_KEY, Value: []byte(event.RecordKey)})
	}
	str := fmt.Sprint(event.RetryCount)
	out = append(out, kafka.Header{Key: H_RETRY_COUNT, Value: []byte(str)})
	out = append(out, kafka.Header{Key: H_EVENT_HASH, Value: []byte(event.EventHash)})
	return out
}

func parseMessage(msg *kafka.Message) *Event {
	out := Event{
		message:        msg.Value,
		topicPartition: msg.TopicPartition,
	}
	h := msg.Headers
	out.Topic = *msg.TopicPartition.Topic
	for _, v := range h {
		if v.Key == H_SENDER_APP {
			out.Sender = string(v.Value)
		}
		if v.Key == H_CONTENT_TYPE {
			out.ContentType = string(v.Value)
		}
		if v.Key == H_SENDER_TOPIC {
			out.SenderTopic = string(v.Value)
		}
		if v.Key == H_RETRY_COUNT {
			retry, _ := strconv.Atoi(string(v.Value))
			if retry > retryCounter {
				retry = retryCounter
			}
			out.RetryCount = retry
		}
		if v.Key == H_EVENT_HASH {
			out.EventHash = string(v.Value)
		}
		if v.Key == H_RECORD_KEY {
			out.RecordKey = string(v.Value)
		}
	}
	return &out
}

func (e *Event) Unmarshal(p interface{}) error {
	if e.ContentType != H_CT_JSON {
		return fmt.Errorf("unkown content type: [%v]", e.ContentType)
	}
	err := json.Unmarshal(e.message, p)
	return err
}

package eventbus

import (
	"github.com/sirupsen/logrus"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/ggrrrr/notification-system/common-lib/config"
)

const (
	CFG_RETRY_TOPIC   = "eventbus.retry.topic"
	CFG_DONE_TOPIC    = "eventbus.done.topic"
	CFG_ERROR_TOPIC   = "eventbus.error.topic"
	CFG_RETRY_COUNTER = "eventbus.error.counter"

	H_TOPIC        = "topic"
	H_SENDER_APP   = "sender_app"
	H_SENDER_TOPIC = "sender_topic"
	H_RETRY_COUNT  = "retry_counter"
	H_EVENT_HASH   = "event_hash"
	H_RECORD_KEY   = "record_key"
	H_CONTENT_TYPE = "content_type"

	H_CT_JSON = "application/json"
)

var commonConfigs []string = []string{
	"bootstrap.servers",
	"sasl.mechanisms",
	"security.protocol",
	"sasl.username",
	"sasl.password",
	"client.id",
	"group.instance.id",
}

var consumerConfigs []string = []string{
	"group.id",
	"auto.offset.reset",
	"client.id",
	"group.instance.id",
	"go.events.channel.enable",
	"go.application.rebalance.enable",
	"enable.auto.commit",
}

var (
	kafkaConfig  kafka.ConfigMap
	retryTopic   string
	doneTopic    string
	errorTopic   string
	retryCounter int
)

func init() {
	retryTopic = config.GetString("", CFG_RETRY_TOPIC)
	doneTopic = config.GetString("", CFG_DONE_TOPIC)
	errorTopic = config.GetString("", CFG_ERROR_TOPIC)
	retryCounter = config.GetInt("", CFG_RETRY_COUNTER, 5)

	kafkaConfig = kafka.ConfigMap{}
	for key := range commonConfigs {
		setConfigValue(kafkaConfig, commonConfigs[key], "")
	}
	logrus.WithFields(logrus.Fields{
		"retryTopic":   retryTopic,
		"doneTopic":    doneTopic,
		"errorTopic":   errorTopic,
		"retryCounter": retryCounter,
	}).Info("init.")
}

func GetRetryTopic() string {
	return retryTopic
}

func GetDoneTopic() string {
	return doneTopic
}

func GetRrrorTopic() string {
	return errorTopic
}

func GetRetryCounter() int {
	return retryCounter
}

func hasString(val kafka.ConfigValue) bool {
	if config.ItoS(val) == "" {
		return false
	}
	return true
}

func consumerConfig(prefix string) kafka.ConfigMap {
	cfg := kafka.ConfigMap{}
	for key := range consumerConfigs {
		setConfigValue(cfg, consumerConfigs[key], prefix)
	}
	return cfg
}

// read config from viper
func setConfigValue(target kafka.ConfigMap, key, prefix string) {
	str := config.GetString(prefix, key)
	if str != "" {
		logrus.Infof("param:%v -> %v = %v", prefix, key, str)
		target[key] = str
	}
}

// Make new copy of common config
func copyConfig(org, inst kafka.ConfigMap) kafka.ConfigMap {
	out := kafka.ConfigMap{}
	for k, v := range org {
		out[k] = v
	}
	for k, v := range inst {
		out[k] = v
	}
	return out
}

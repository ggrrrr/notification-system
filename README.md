# Notification app API

## Description

The purpose of this system is to enable push notifications to third party actors/systems via preconfigured providers ( channels).
Each notification needs to match to a channel. This process is in two steps:

1. notification will have a `channel` field which needs to match to the destination channel
2. Each `channel` has it own additional parameters which needs to be validated.

If the above two step are verified the channel subsystem will attempt to push the notification to the provider.
In case of error, it will forward it to the `RetryQueue` system for later attempt.

## Api Interface

The system is using EventBus interface, based on Apache Kafka Pub/Sub service.
To send notificaiton message you need to `publish` a message to a Kafka `topic`
with the following details:

- `topic name`: `notification_request`
- Headers
  - `sender_app`: name of the system sender
  - `sender_topic`: topic on which a result of the operationc
  - `content_type`: _required_ payload encoidng: currently only `application/json` is supported
  -
- payload: JSON with the following feilds:
  - channel `string`: name of the destination chennal
  - payload `map[string]any`: channel specific fields

Example payloads are listed bellow.

Please see [commont-lib/eventbus/events_tests.go](commont-lib/eventbus/events_tests.gomd) for examples.

## Available Channels

### Dummy Channel

Just a basic channel, not for production.

### SMS Channel

Mobile phone SMS channel, provider is www.infobip.com

Config:

```yaml
sms:
  enable: true
  from: InfoSMS
  api:
  	token: <TOKEN>
	url: https://someurl

```

Example notification request payload:

```go
	msg := notification.NotificationData{
		Channel: "slack",
		Body: map[string]interface{}{
			"text":  "text me",
			"from":  "InfoSMS",
			"phone": "12312313",
		},
	}

```

API from [https://www.infobip.com/docs/api#channels/sms/send-sms-message]

### Slack Channel

Slack is very popular chat system. www.slack.com

Config:

```yaml
slack:
  enable: true
  api:
  	token: <TOKEN>
	url: https://slack.com/api/chat.postMessage

```

Example notification request payload:

```go
	msg := notification.NotificationData{
		Channel: "slack",
		Body: map[string]interface{}{
			"text":    "text me",
			"channel": "asdasd",
		},
	}

```

API from [https://api.slack.com/methods/chat.postMessage]

### Email channel

Basic email notification.

Config:

```yaml
email:
  enable: true
  smtp:
  	server: smtp.example.com
	port: 25
	username: john@example.com
	password: <PASS>
	sender: john@example.com
	tls: false

```

Example notification request payload:

```go
	msg := notification.NotificationData{
		Channel: "email",
		Body: map[string]interface{}{
			"text":    "text me",
			"to":      "asdasd",
			"subject": "new subject",
		},
	}

```

API from [https://pkg.go.dev/net/smtp]

## Configuration

Configuration fileA [configs/app.yaml](configs/app.yaml) will be loaded. All variable can be overwritten with OS environment variables.
Where each YAML path will be Upper case with `_` delimiter example:

```yaml
dummy:
	enable: true
```

will be

```
export DUMMY_ENABLE=true
```

## Adding new channel

1. You can copy [notifications-app/channels/dummy](notifications-app/channels/dummy) to `notifications-app/channels/mychannel`
2. Modify `New()` method to match your channel configuration needs.
3. Update `Push` interface method based on your channel provider API
4. In `main.go` add the following lines:

   ```go
   	mychannelC, err := mychannel.New()
   	if err == nil {
   		registry.Register(mychannelC)
   	}

   ```

# Usefull tools

## Publish basic kafka message

```
kafka-console-producer --topic notification_request --bootstrap-server localhost:9092
```

## Consumne all kafka messages from topic

```
kafka-console-consumer --topic retry_topic --bootstrap-server localhost:9092 --from-beginning --property print.headers=true
```

## Debug inject notification with cURL

```
curl -H "Content-Type: application/json" \
	-X POST \
	-d '{"channel":"dummy","body":{"dummy":"retry"}}' \
	localhost:6021/process
```

## Testing and local development

### Dependencies

`Docker`
`docker-compose`

from [local/docker-compose.yaml](local/docker-compose.yaml) you can install kafka on local PC by running

```
docker-compose up -d zookeeper # zookeeper needed
sleep 10 # this will allow zookeepr to start
docker-compose up -d broker  # kafka broker

```

Load env by running `source env-local.sh`, this will export all variable from:

`.env.local`

Example file:

```bash
cat > .env.local <<EOF
LISTEN_ADDR=:6021

BOOTSTRAP_SERVERS=localhost:9092

EVENT_TOPIC=notification_request
EVENT_GROUP_ID=mygroup
EVENT_AUTO_OFFSET_RESET=beginning

EVENTBUS_RETRY_TOPIC=retry_topic
EVENTBUS_ERROR_TOPIC=error_topic
EVENTBUS_DONE_TOPIC=done_topic

RETRY_QUEUE_GROUP_ID=local
RETRY_QUEUE_AUTO_OFFSET_RESET=beginning

TEST_GROUP_ID=local12331
TEST_GROUP_INSTANCE_ID=static123
TEST_AUTO_OFFSET_RESET=smallest
TEST_CLIENT_ID=test_app

RETRY_QUEUE_SLEEP_TIME=2s

SLACK_ENABLE=true
SLACK_API_URL=https://slack.com/api/chat.postMessage

SMS_ENABLE=true
SMS_API_TOKEN=SOMETOKEN
SMS_API_URL=
EOF
```

##### Run application from console:

###### Run retry-queue-svc from console:

```bash
source env-local.sh
HTTP_LISTEN_ADDR=:8081 go run retry-queue-svc/main.go
```

###### Run notifications-appfrom console:

```bash
source env-local.sh
HTTP_LISTEN_ADDR=:8081 go run notifications-app/main.go

```

##### Production deployment

You can use:

- [notifications-app.yaml](k8s/notifications-app.yaml)
- [retry-queue-svc.yaml](k8s/retry-queue-svc.yaml)

as base to create kubernetes deployment.

:beers:

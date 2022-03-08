package notification

const (
	C_ENABLE = "enable"
)

type EVENT_RESULT int

const (
	EVENT_DONE  EVENT_RESULT = 0
	EVENT_RETRY              = 1
	EVENT_ERROR              = 2
)

// Notification event
type NotificationData struct {
	Channel string                 `json:"channel"`
	Body    map[string]interface{} `json:"body"`
}

type NotificationService interface {
	Push(msg *NotificationData) (EVENT_RESULT, error)
	Name() string
}

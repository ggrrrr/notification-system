package notification_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ggrrrr/notification-system/notifications-app/notification"
)

type testT struct {
	json   string
	result *notification.NotificationData
	err    error
}

var (
	tests = []testT{
		{
			json:   `{"channel":"dummy","body":{"text":"text"}}`,
			result: &notification.NotificationData{Channel: "dummy", Body: map[string]interface{}{"text": "text"}},
			err:    nil,
		},
		{
			json:   `{"channel":"dummy","body":{"text":"text"}}`,
			result: &notification.NotificationData{Channel: "dummy", Body: map[string]interface{}{"text": "text"}},
			err:    nil,
		},
		{
			json:   `{"channel":"slack","body":{"text":"text"}}`,
			result: &notification.NotificationData{Channel: "slack", Body: map[string]interface{}{"text": "text"}},
			err:    nil,
		},
		{
			json:   `{"channel":"slack","body":{"text":"text","text1":"text12"}}`,
			result: &notification.NotificationData{Channel: "slack", Body: map[string]interface{}{"text": "text", "text1": "text12"}},
			err:    nil,
		},
	}
)

func verify(t *testing.T, v testT, r *notification.NotificationData, err error) {
	if !reflect.DeepEqual(v.result, r) {
		t.Errorf(" dont match %+v", v.result)
		t.Errorf(" dont match %+v", r)
	}
}

func Test1(t *testing.T) {
	for _, v := range tests {
		var r notification.NotificationData
		err := json.Unmarshal([]byte(v.json), &r)
		verify(t, v, &r, err)
	}
}

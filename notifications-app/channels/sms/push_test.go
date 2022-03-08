package sms_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ggrrrr/notification-system/notifications-app/channels/sms"
	"github.com/ggrrrr/notification-system/notifications-app/notification"
)

func TestOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"bulkId": "2034072219640523072",
			"messages": [
			  {
				"to": "41793026727",
				"status": {
				  "groupId": 1,
				  "groupName": "PENDING",
				  "id": 26,
				  "name": "PENDING_ACCEPTED",
				  "description": "Message sent to next instance"
				},
				"messageId": "2250be2d4219-3af1-78856-aabe-1362af1edfd2"
			  }
			]
		  }`)
	}))
	defer ts.Close()
	t.Setenv("SMS_API_URL", ts.URL)
	sms, err := sms.New()
	if err != nil {
		t.Fatalf("%v", err)
	}

	msg := notification.NotificationData{
		Channel: "slack",
		Body: map[string]interface{}{
			"text":  "text me",
			"from":  "InfoSMS",
			"phone": "12312313",
		},
	}

	ok, err := sms.Push(&msg)

	if ok != notification.EVENT_DONE || err != nil {
		t.Errorf("ok %v %v", ok, err)
	}

}

func TestErr(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"requestError": {
			  "serviceException": {
				"messageId": "string",
				"text": "string"
			  }
			}
		  }`)
	}))
	defer ts.Close()
	t.Setenv("SMS_API_URL", ts.URL)
	sms, err := sms.New()
	if err != nil {
		t.Fatalf("%v", err)
	}

	msg := notification.NotificationData{
		Channel: "slack",
		Body:    map[string]interface{}{"text": "text me", "from": "InfoSMS", "phone": "12312313"},
	}

	ok, err := sms.Push(&msg)
	if ok != notification.EVENT_RETRY {
		t.Errorf("ok %v %v", ok, err)
	}

}

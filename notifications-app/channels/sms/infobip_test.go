package sms_test

import (
	"encoding/json"
	"testing"

	"github.com/ggrrrr/notification-system/notifications-app/channels/sms"
)

func TestXxx(t *testing.T) {
	sms1Json1 := `{"messages":[{"from":"InfoSMS","destinations":[{"to":"41793026727"}],"text":"This is a sample message"}]}`
	sms1 := sms.Request{
		Messages: []sms.Message{
			{
				From: "InfoSMS",
				Text: "This is a sample message",
				Dest: []sms.To{{To: "41793026727"}},
			},
		},
	}
	sms1Json, _ := json.Marshal(sms1)
	if string(sms1Json) != sms1Json1 {
		t.Log(string(sms1Json1))
		t.Log(string(sms1Json))
		t.Errorf("request wrong")
	}

	respJsonOK := `{
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
	  }`
	var responseOK sms.Response

	err := json.Unmarshal([]byte(respJsonOK), &responseOK)
	if err != nil {
		t.Error(err)
	}
	if responseOK.BulkId == "" {
		t.Error("builk not here")

	}

	responseJsonErr := `{
		"requestError": {
		  "serviceException": {
			"messageId": "string",
			"text": "string"
		  }
		}
	  }`
	var responseError sms.Response

	err = json.Unmarshal([]byte(responseJsonErr), &responseError)
	if err != nil {
		t.Error(err)
	}
	if responseError.BulkId != "" {
		t.Error("builk not here")

	}
	if responseError.Error.Exception.MessageId == "" {
		t.Error("MessageId not here")

	}
	t.Logf("ok : %+v", responseOK)
	t.Logf("err: %+v", responseError)

}

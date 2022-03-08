package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ggrrrr/notification-system/notifications-app/notification"
	"github.com/sirupsen/logrus"
)

func (c *cfg) Push(msg *notification.NotificationData) (notification.EVENT_RESULT, error) {
	logrus.Infof("msg: %v", msg)

	phone, ok := msg.Body["phone"]
	if !ok {
		return notification.EVENT_ERROR, fmt.Errorf("missing phone")
	}
	text, ok := msg.Body["text"]
	if !ok {
		return notification.EVENT_ERROR, fmt.Errorf("missing text")
	}
	from, ok := msg.Body["from"]
	if !ok {
		from = c.from
	}

	url := fmt.Sprintf(c.url)
	infoBip := Request{
		Messages: []Message{
			{
				From: fmt.Sprintf("%v", from),
				Text: fmt.Sprintf("%v", text),
				Dest: []To{{To: fmt.Sprintf("%v", phone)}},
			},
		},
	}

	body, err := json.Marshal(infoBip)
	if err != nil {
		return notification.EVENT_ERROR, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return notification.EVENT_ERROR, err
	}
	request.Header.Add("Authorization", c.token)
	request.Header.Add("Content-type", "application/json")
	request.Header.Add("Accept", "application/json")

	client := &http.Client{}
	logrus.Infof("url: %s", url)
	// logrus.Infof("token: %s", bearer)
	// logrus.Debug("body: ", string(body))
	res, err := client.Do(request)
	if err != nil {
		logrus.Errorf("response: %v", res.Status)
		logrus.Error(err)
		return notification.EVENT_RETRY, err
	}
	logrus.Infof("response: %v", res.Status)
	body1, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return notification.EVENT_RETRY, err
	}
	var response Response
	err = json.Unmarshal(body1, &response)
	if err != nil {
		return notification.EVENT_RETRY, err
	}
	logrus.Debugf("response: %v", response)
	if response.BulkId == "" {
		return notification.EVENT_RETRY, fmt.Errorf("infobip error: %+v", response.Error)
	}
	return notification.EVENT_DONE, nil
}

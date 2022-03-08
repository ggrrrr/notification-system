package slack

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
	url := fmt.Sprintf(c.url)
	_, ok := msg.Body["text"]
	if !ok {
		return notification.EVENT_ERROR, fmt.Errorf("missing text")
	}
	_, ok = msg.Body["channel"]
	if !ok {
		return notification.EVENT_ERROR, fmt.Errorf("missing channel")
	}

	body, err := json.Marshal(msg.Body)
	if err != nil {
		return notification.EVENT_ERROR, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return notification.EVENT_ERROR, err
	}
	bearer := fmt.Sprintf("Bearer %s", c.token)
	request.Header.Add("Authorization", bearer)
	request.Header.Add("Content-type", "application/json")

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
	var slackResponse SlackResponse
	err = json.Unmarshal(body1, &slackResponse)
	if err != nil {
		return notification.EVENT_RETRY, err
	}
	if !slackResponse.Ok {
		return notification.EVENT_RETRY, fmt.Errorf("slack error: %v", slackResponse.Error)

	}
	return notification.EVENT_DONE, nil
}

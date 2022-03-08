package eventproc

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ggrrrr/notification-system/common-lib/app"
	"github.com/ggrrrr/notification-system/common-lib/eventbus"
	"github.com/ggrrrr/notification-system/notifications-app/notification"
	"github.com/sirupsen/logrus"
)

func HttpHandle(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get(string(app.HTTP_CONTENT_TYPE))
	if contentType != eventbus.H_CT_JSON {
		logrus.Errorf("unsupported content type")
		app.HttpError(w, http.StatusBadRequest,
			"content-type",
			fmt.Errorf("bad content-type: %v", contentType),
		)
		return
	}
	logrus.Infof("%v", r.URL)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("read body error: %+v", err)
		app.HttpError(w, http.StatusBadRequest,
			"error ioutil",
			err,
		)
		return
	}

	var msg notification.NotificationData
	logrus.Debugf("msg: %+v", msg)
	event, err := eventbus.NewFromJson(topic, body, &msg)
	if err != nil {
		app.HttpError(w, http.StatusBadRequest,
			"json",
			err,
		)
		return
	}
	// logrus.Debugf("%v", event)
	go EventHandler(event)
	app.HttpOk(w, "accepted")
}

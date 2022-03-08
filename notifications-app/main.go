package main

import (
	"os"
	"os/signal"

	"github.com/ggrrrr/notification-system/common-lib/app"
	"github.com/ggrrrr/notification-system/notifications-app/channels/dummy"
	"github.com/ggrrrr/notification-system/notifications-app/channels/email"
	"github.com/ggrrrr/notification-system/notifications-app/channels/slack"
	"github.com/ggrrrr/notification-system/notifications-app/channels/sms"
	"github.com/ggrrrr/notification-system/notifications-app/eventproc"
	"github.com/ggrrrr/notification-system/notifications-app/notification"
	"github.com/sirupsen/logrus"
)

var (
	osSignals = make(chan os.Signal, 1)
)

func main() {
	app.Configure("notification-app", osSignals)
	app.HandleFunc("/process", eventproc.HttpHandle)

	notification.Init()

	dummy, err := dummy.New()
	if err == nil {
		notification.Add(dummy)
	}

	smsC, err := sms.New()
	if err == nil {
		notification.Add(smsC)
	}

	slackC, err := slack.New()
	if err == nil {
		notification.Add(slackC)
	}
	emailC, err := email.New()
	if err == nil {
		notification.Add(emailC)
	}

	err = eventproc.Configure()
	if err != nil {
		panic(err)
	}
	go eventproc.EventLoop()
	eventproc.Start()
	go app.Start()
	signal.Notify(osSignals, os.Interrupt)
	logrus.Printf("os.signal: %v", <-osSignals)
	logrus.Printf("end.")
}

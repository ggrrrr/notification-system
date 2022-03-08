package main

import (
	"os"
	"os/signal"

	"github.com/ggrrrr/notification-system/common-lib/app"
	"github.com/ggrrrr/notification-system/retry-queue-svc/eventproc"
	"github.com/sirupsen/logrus"
)

var (
	err       error
	osSignals = make(chan os.Signal, 1)
)

func main() {
	app.Configure("retry-queue-app", osSignals)

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

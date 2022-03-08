package eventbus

type EventProducer interface {
	PublishEvent(event *Event)
	PublishRetry(event *Event)
	PublishDone(event *Event)
	PublishError(event *Event)
	// Close()
	WaitAll()
}

type EventConsumer interface {
	Start(eventChan chan *Event) error
	Pause(event *Event)
	Resume(event *Event)
	Close()
}

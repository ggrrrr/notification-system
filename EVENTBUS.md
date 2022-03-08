# EventBus APIs

## Purspose

Asynchronous way of communication between decoupled systems - microservices. Usually such communication can be group per type:

- update requests: for example update on account balances, create new records or other operations after which the system will have different state.
  These are request for chage in the state.
- Notifications: a records was created, balance was updated, errors etc.
  These are notification after change in the state.

To implement Producer of events you can use the folloing code:

```
// Create producer instance in your main.go
// this can be only one per application
producer, err := eventbus.NewProducer(nil)

// To create event use the following
// Create event payload
payload := struct{ Key string }{Key: uuid.New().String()}

// Create to event message
testEvent1 := eventbus.NewEvent(topic, payload)



// configure a topic

```

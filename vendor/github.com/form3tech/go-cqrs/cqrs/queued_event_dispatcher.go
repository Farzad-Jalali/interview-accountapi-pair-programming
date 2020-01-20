package cqrs

import (
	"sync"

	"reflect"

	"context"

	"encoding/json"

	"github.com/form3tech/go-messaging/messaging"
)

var oneQueuedEventDispatcher sync.Once
var queuedEventDispatcherInstance *QueuedEventDispatcher

type QueuedEventDispatcher struct {
	receiver           messaging.Receiver
	sender             messaging.Sender
	QueueNamer         QueueNamingStrategy
	subscribedQueues   map[string]interface{}
	eventTypes         map[string]reflect.Type
	internalDispatcher *inMemoryEventDispatcher
}

func GetQueuedEventDispatcher(applicationName string, r messaging.Receiver, s messaging.Sender) *QueuedEventDispatcher {
	oneQueuedEventDispatcher.Do(func() {
		queuedEventDispatcherInstance = newQueuedEventDispatcher(applicationName, r, s)
	})

	return queuedEventDispatcherInstance
}

func newQueuedEventDispatcher(applicationName string, r messaging.Receiver, s messaging.Sender) *QueuedEventDispatcher {
	return &QueuedEventDispatcher{
		receiver:   r,
		sender:     s,
		QueueNamer: PerTypeQueueNamingStrategy{ApplicationName: applicationName},
		internalDispatcher: &inMemoryEventDispatcher{
			handlers: make(map[string][]eventHandler),
		},
		subscribedQueues: make(map[string]interface{}),
		eventTypes:       make(map[string]reflect.Type),
	}
}

func (a *QueuedEventDispatcher) Dispatch(ctx *context.Context, e interface{}) error {
	eventType := reflect.TypeOf(e).String()
	return a.sender.Send(a.QueueNamer.GetQueueNameFor(eventType), messaging.CreateMessage(e))
}

func (a *QueuedEventDispatcher) handleMessage(ctx context.Context, message messaging.Message) error {
	eventType, typeInfoExists := message.MessageAttributes[messaging.TypeInfoAttribute]
	if !typeInfoExists {
		return nil
	}

	typeInfo, canHandle := a.eventTypes[eventType.(string)]
	if !canHandle {
		return nil
	}

	event := reflect.New(typeInfo)
	if err := json.Unmarshal([]byte(message.Body.(string)), event.Interface()); err != nil {
		return err
	}

	ctx = context.WithValue(ctx, applicationContext, true)
	return a.internalDispatcher.Dispatch(&ctx, event.Elem().Interface())
}

func (a *QueuedEventDispatcher) RegisterEventHandler(handler interface{}) error {
	if err := a.internalDispatcher.RegisterEventHandler(handler); err != nil {
		return err
	}

	eventType, _, err := a.internalDispatcher.getEventFromHandler(handler)
	if err != nil {
		return err
	}

	queueName := a.QueueNamer.GetQueueNameFor(eventType.String())

	if _, subscribed := a.subscribedQueues[queueName]; !subscribed {
		err := a.receiver.SubscribeMessage(queueName, a.handleMessage)
		if err != nil {
			return err
		}
		a.subscribedQueues[queueName] = a.handleMessage
	}
	a.eventTypes[eventType.String()] = eventType

	return nil
}

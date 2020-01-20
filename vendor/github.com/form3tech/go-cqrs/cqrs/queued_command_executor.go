package cqrs

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/form3tech/go-messaging/messaging"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var oneQueuedExecutor sync.Once
var queuedInstance *QueuedCommandExecutor

const applicationContext string = "security.application_context"

type QueuedCommandExecutor struct {
	receiver         messaging.Receiver
	sender           messaging.Sender
	QueueNamer       QueueNamingStrategy
	subscribedQueues map[string]interface{}
	commandTypes     map[string]reflect.Type
	internalExecutor *inMemoryCommandExecutor
}

func GetQueuedCommandExecutor(db *sqlx.DB, applicationName string, receiver messaging.Receiver, sender messaging.Sender) *QueuedCommandExecutor {
	oneQueuedExecutor.Do(func() {
		queuedInstance = newQueuedCommandExecutor(db, applicationName, receiver, sender)
	})

	return queuedInstance
}

func newQueuedCommandExecutor(db *sqlx.DB, applicationName string, r messaging.Receiver, s messaging.Sender) *QueuedCommandExecutor {
	return &QueuedCommandExecutor{
		receiver:   r,
		sender:     s,
		QueueNamer: PerTypeQueueNamingStrategy{ApplicationName: applicationName},
		internalExecutor: &inMemoryCommandExecutor{
			handlers: make(map[string]h),
			db:       db,
		},
		subscribedQueues: make(map[string]interface{}),
		commandTypes:     make(map[string]reflect.Type),
	}
}

func (d *QueuedCommandExecutor) Execute(ctx *context.Context, organisationId *uuid.UUID, c interface{}) error {
	commandType := reflect.TypeOf(c).String()

	_, ok := d.internalExecutor.handlers[commandType]
	if !ok {
		return fmt.Errorf("no command handler registered of type: %s", commandType)
	}

	queueName := d.QueueNamer.GetQueueNameFor(commandType)
	return d.sender.Send(queueName, messaging.CreateMessage(c))
}

func (a *QueuedCommandExecutor) handleMessage(ctx context.Context, message messaging.Message) error {

	eventType, typeInfoExists := message.MessageAttributes[messaging.TypeInfoAttribute]
	if !typeInfoExists {
		return nil
	}

	typeInfo, canHandle := a.commandTypes[eventType.(string)]
	if !canHandle {
		return nil
	}

	command := reflect.New(typeInfo)
	if err := json.Unmarshal([]byte(message.Body.(string)), command.Interface()); err != nil {
		return err
	}

	ctx = context.WithValue(ctx, applicationContext, true)
	return a.internalExecutor.Execute(&ctx, nil, command.Elem().Interface())
}

func (d *QueuedCommandExecutor) RegisterCommandHandler(handler interface{}, checkPermissionsFn func(*context.Context, *uuid.UUID) error) error {
	if err := d.internalExecutor.RegisterCommandHandler(handler, checkPermissionsFn); err != nil {
		return err
	}

	commandType, _, _, err := d.internalExecutor.getCommandFromHandler(handler)
	if err != nil {
		return err
	}

	queueName := d.QueueNamer.GetQueueNameFor(commandType.String())
	if _, subscribed := d.subscribedQueues[queueName]; !subscribed {
		err := d.receiver.SubscribeMessage(queueName, d.handleMessage)
		if err != nil {
			return err
		}
		d.subscribedQueues[queueName] = d.handleMessage
	}
	d.commandTypes[commandType.String()] = commandType

	return nil
}

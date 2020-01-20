package executors

import (
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/settings"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/log"
	"github.com/form3tech/go-cqrs/cqrs"
	"github.com/form3tech/go-messaging/messaging"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var Receiver messaging.Receiver
var Sender messaging.Sender
var InMemoryCommandExecutor cqrs.CommandExecutor
var QueuedCommandExecutor cqrs.CommandExecutor
var QueuedEventDispatcher cqrs.EventDispatcher
var InMemoryEventDispatcher cqrs.EventDispatcher
var QueryExecutor cqrs.QueryExecutor

func Configure(db *sqlx.DB) {
	messageVisibilityTimeout := viper.GetInt("MessageVisibilityTimeout")

	Sender = messaging.NewSqsSender()
	Receiver = messaging.NewSqsReceiverBuilder().
		WithLogger(log.StandardLogger()).
		WithVisibilityTimeout(int64(messageVisibilityTimeout)).
		Build()

	InMemoryCommandExecutor = cqrs.GetInMemoryCommandExecutor(db)
	QueryExecutor = cqrs.GetQueryExecutor(db)
	InMemoryEventDispatcher = cqrs.GetInMemoryEventDispatcher()

	queuedExecutor := cqrs.GetQueuedCommandExecutor(db, settings.ApiName, Receiver, Sender)
	queuedExecutor.QueueNamer = cqrs.SharedQueueNamingStrategy{ApplicationName: settings.ApiName, QueueName: "commands"}

	queuedDispatcher := cqrs.GetQueuedEventDispatcher(settings.ApiName, Receiver, Sender)
	queuedDispatcher.QueueNamer = cqrs.SharedQueueNamingStrategy{ApplicationName: settings.ApiName, QueueName: "events"}

	QueuedEventDispatcher = queuedDispatcher
	QueuedCommandExecutor = queuedExecutor
}

package eventhandlers

import (
	"github.com/form3tech/go-messaging/messaging"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/executors"
)

const (
	ExternalEventDestinationName = "form3-events"
)

var form3EventSender messaging.Sender

func Configure() {
	var err error
	form3EventSender = messaging.NewSnsSender()
	err = executors.InMemoryEventDispatcher.RegisterEventHandler(Form3EventNotificationEventHandler)
	if err != nil {
		panic(err)
	}

	//TODO: Register event handlers here
}

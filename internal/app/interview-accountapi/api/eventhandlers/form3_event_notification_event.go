package eventhandlers

import (
	"github.com/form3tech/go-messaging/messaging"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/events"
	"github.com/pkg/errors"
)

func Form3EventNotificationEventHandler(e events.Form3EventNotificationEvent) error {
	err := form3EventSender.Send(ExternalEventDestinationName, messaging.Message{
		Body: e.Event,
	})
	if err != nil {
		return errors.Wrapf(err, "unable to send event notification. msg:%+v, err:%+v", e, err)
	}
	return nil
}

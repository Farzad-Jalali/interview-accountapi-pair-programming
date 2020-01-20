package events

import "github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/internalmodels"

type Form3EventNotificationEvent struct {
	Event *internalmodels.Form3Event
}

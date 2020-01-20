package cqrs

import (
	"fmt"
	"strings"
)

type QueueNamingStrategy interface {
	GetQueueNameFor(typeName string) string
}

type PerTypeQueueNamingStrategy struct {
	ApplicationName string
}

type SharedQueueNamingStrategy struct {
	ApplicationName string
	QueueName       string
}

func (s PerTypeQueueNamingStrategy) GetQueueNameFor(typeName string) string {
	result := strings.Replace(strings.ToLower(typeName), "api.", "", 1)
	result = strings.Replace(result, "events.", "", 1)
	result = strings.Replace(result, "event", "", 1)
	return fmt.Sprintf("%s-%s", s.ApplicationName, result)
}

func (s SharedQueueNamingStrategy) GetQueueNameFor(typeName string) string {
	return fmt.Sprintf("%s-%s", s.ApplicationName, s.QueueName)
}

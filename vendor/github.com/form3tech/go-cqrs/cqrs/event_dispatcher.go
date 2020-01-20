package cqrs

import (
	"context"
)

type EventDispatcher interface {
	Dispatch(ctx *context.Context, e interface{}) error
	RegisterEventHandler(handler interface{}) error
}

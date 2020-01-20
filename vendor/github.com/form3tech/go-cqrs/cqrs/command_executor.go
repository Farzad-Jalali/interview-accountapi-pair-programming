package cqrs

import (
	"context"

	"github.com/google/uuid"
)

type CommandExecutor interface {
	Execute(ctx *context.Context, organisationId *uuid.UUID, command interface{}) error
	RegisterCommandHandler(handler interface{}, checkPermissionsFn func(ctx *context.Context, organisationId *uuid.UUID) error) error
}

type h struct {
	handler            interface{}
	checkPermissionsFn func(ctx *context.Context, organisationId *uuid.UUID) error
	usesDb             bool
	usesContext        bool
}

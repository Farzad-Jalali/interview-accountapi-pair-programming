package cqrs

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

var oneExecutor sync.Once
var inMemoryInstance *inMemoryCommandExecutor

type inMemoryCommandExecutor struct {
	handlers map[string]h
	db       *sqlx.DB
}

func GetInMemoryCommandExecutor(db *sqlx.DB) CommandExecutor {
	oneExecutor.Do(func() {
		inMemoryInstance = &inMemoryCommandExecutor{
			handlers: make(map[string]h),
			db:       db,
		}
	})

	return inMemoryInstance
}

func (d *inMemoryCommandExecutor) Execute(ctx *context.Context, organisationId *uuid.UUID, c interface{}) (err error) {

	commandType := reflect.TypeOf(c)

	h, ok := d.handlers[commandType.String()]

	if !ok {
		return fmt.Errorf("no command handler registered of type: %s", commandType)
	}

	err = h.checkPermissionsFn(ctx, organisationId)

	if err != nil {
		return
	}

	handler := h.handler
	method := reflect.ValueOf(handler)

	var result []reflect.Value

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("handler for '%s' encountered a panic: %s", commandType, e)
		}
	}()

	if ctx != nil {
		annotatedContext := context.WithValue(*ctx, "handler", runtime.FuncForPC(method.Pointer()).Name())
		ctx = &annotatedContext
	}

	if h.usesDb && !h.usesContext {
		result = method.Call([]reflect.Value{reflect.ValueOf(d.db), reflect.ValueOf(c)})
	} else if !h.usesDb && h.usesContext {
		result = method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(c)})
	} else if h.usesDb && h.usesContext {
		result = method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(d.db), reflect.ValueOf(c)})
	} else {
		result = method.Call([]reflect.Value{reflect.ValueOf(c)})
	}

	if result[0].Interface() != nil {
		err, ok = result[0].Interface().(error)
		if !ok {
			err = fmt.Errorf("handler for '%s' returned non-error interface as first return value", commandType)
		} else {
			log.WithField("handler", runtime.FuncForPC(method.Pointer()).Name()).Errorf("handler for '%s' failed: %+v", commandType, err)
		}
	}
	return
}

func (d *inMemoryCommandExecutor) RegisterCommandHandler(handler interface{}, checkPermissionsFn func(*context.Context, *uuid.UUID) error) error {
	commandType, usesDb, usesContext, err := d.getCommandFromHandler(handler)

	if err != nil {
		return err
	}

	if _, ok := d.handlers[commandType.String()]; ok {
		return fmt.Errorf("handler already registered for %s only one handler allowed per command", commandType)
	}

	d.handlers[commandType.String()] = h{
		handler:            handler,
		checkPermissionsFn: checkPermissionsFn,
		usesDb:             usesDb,
		usesContext:        usesContext,
	}

	return nil
}

func (*inMemoryCommandExecutor) getCommandFromHandler(handler interface{}) (commandType reflect.Type, usesDb bool, usesContext bool, err error) {
	var dbType *sqlx.DB
	var ctxType *context.Context
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		err = fmt.Errorf("%s is not a func must be a func with one arg", handlerType.Kind())
		return
	}
	if handlerType.NumIn() < 1 || handlerType.NumIn() > 3 {
		err = fmt.Errorf("%s must only have one to three arguments", handlerType.Kind())
		return
	}

	if handlerType.NumIn() == 1 {
		commandType = handlerType.In(0)
	}

	if handlerType.NumIn() == 2 {
		if handlerType.In(0) == reflect.TypeOf(dbType) {
			usesDb = true
		} else if handlerType.In(0) == reflect.TypeOf(ctxType) {
			usesContext = true
		} else {
			err = fmt.Errorf("when using 2 arguments first argument should be of type *sqlx.db or *context.Context")
			return
		}

		commandType = handlerType.In(1)
	}

	if handlerType.NumIn() == 3 {
		if handlerType.In(0) == reflect.TypeOf(ctxType) {
			usesContext = true
		}

		if handlerType.In(1) == reflect.TypeOf(dbType) {
			usesDb = true
		}

		if !usesContext {
			err = fmt.Errorf("when using 3 arguments first argument should be of type *context.Context")
			return
		}

		if !usesDb {
			err = fmt.Errorf("when using 3 arguments second argument should be of type *sqlx.DB")
			return
		}
		commandType = handlerType.In(2)
	}
	return
}

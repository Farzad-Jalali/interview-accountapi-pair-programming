package cqrs

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var oneInMemoryEventDispatcher sync.Once
var inMemoryEventDispatcherInstance *inMemoryEventDispatcher

type eventHandler struct {
	handler     interface{}
	usesContext bool
}

type inMemoryEventDispatcher struct {
	handlers map[string][]eventHandler
}

type eventError struct {
	message string
}

func (e *eventError) Error() string {
	return e.message
}

func newEventError(errors []error) error {

	if len(errors) == 0 {
		return nil
	}

	var errorMessages []string

	for _, m := range errors {
		errorMessages = append(errorMessages, m.Error())
	}

	return &eventError{
		message: strings.Join(errorMessages, ","),
	}

}

func GetInMemoryEventDispatcher() EventDispatcher {
	oneInMemoryEventDispatcher.Do(func() {
		inMemoryEventDispatcherInstance = &inMemoryEventDispatcher{
			handlers: make(map[string][]eventHandler),
		}
	})

	return inMemoryEventDispatcherInstance
}

func (d *inMemoryEventDispatcher) Dispatch(ctx *context.Context, e interface{}) error {

	eventType := reflect.TypeOf(e)

	s := eventType.String()
	handlers, ok := d.handlers[s]

	if !ok {
		return fmt.Errorf("no handler registered for type '%s'", eventType)
	}

	eventValue := reflect.ValueOf(e)
	var errors []error

	var runFunc = func(handler eventHandler) (err error) {

		method := reflect.ValueOf(handler.handler)

		var result []reflect.Value

		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("encountered a panic: %s", e)
			}
		}()

		if handler.usesContext {
			if ctx != nil {
				annotatedContext := context.WithValue(*ctx, "handler", runtime.FuncForPC(method.Pointer()).Name())
				ctx = &annotatedContext
			}
			result = method.Call([]reflect.Value{reflect.ValueOf(ctx), eventValue})
		} else {
			result = method.Call([]reflect.Value{eventValue})
		}

		if result[0].Interface() != nil {
			err, ok = result[0].Interface().(error)
			if !ok {
				err = fmt.Errorf("handler for '%s' returned non-error interface as first return value", eventType)
			} else {
				log.WithField("handler", runtime.FuncForPC(method.Pointer()).Name()).Errorf("handler for '%s' failed: %+v", eventType, err)
			}
		}
		return
	}

	for i, handler := range handlers {
		if err := runFunc(handler); err != nil {
			errors = append(errors, fmt.Errorf("error in handler #%d for '%s': %s", i, eventType, err))
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return newEventError(errors)
}

func (d *inMemoryEventDispatcher) RegisterEventHandler(handler interface{}) error {
	eventType, usesContext, err := d.getEventFromHandler(handler)
	if err != nil {
		return err
	}

	handlers := d.handlers[eventType.String()]

	d.handlers[eventType.String()] = append(handlers, eventHandler{
		handler:     handler,
		usesContext: usesContext,
	})

	return nil
}

func (inMemoryEventDispatcher) getEventFromHandler(handler interface{}) (eventType reflect.Type, usesContext bool, err error) {
	handlerType := reflect.TypeOf(handler)

	if handlerType.Kind() != reflect.Func || handlerType.NumIn() < 1 || handlerType.NumIn() > 2 || handlerType.NumOut() != 1 {
		err = fmt.Errorf("handler must be a func with one argument and return error")
		return
	}

	if handlerType.NumIn() == 2 {
		var ctxType *context.Context
		if handlerType.In(0) != reflect.TypeOf(ctxType) {
			err = fmt.Errorf("when using 2 arguments first argument should be of type *context.Context")
			return
		}

		usesContext = true
		eventType = handlerType.In(1)
	}

	if handlerType.NumIn() == 1 {
		eventType = handlerType.In(0)
	}
	return
}

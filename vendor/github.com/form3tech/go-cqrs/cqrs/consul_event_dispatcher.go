package cqrs

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/consul/api"
)

var consulDispatcher sync.Once
var consulDispatcherInstance *consulEventDispatcher

type consulEventDispatcher struct {
	client  *api.Client
	waiters []*consulEventWaiter
}

// Special purpose dispatcher when you want to send an event only to other nodes
// and not using any infrastructure like queue
func GetConsulEventDispatcher() (EventDispatcher, error) {
	var err error
	consulDispatcher.Do(func() {
		consulDispatcherInstance, err = newConsulEventDispatcher()
	})
	return consulDispatcherInstance, err
}

func newConsulEventDispatcher() (*consulEventDispatcher, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return &consulEventDispatcher{
		client:  client,
		waiters: make([]*consulEventWaiter, 0),
	}, nil
}

func (m *consulEventDispatcher) Dispatch(ctx *context.Context, e interface{}) error {

	eventType := reflect.TypeOf(e).String()
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, _, err = m.client.Event().Fire(&api.UserEvent{Name: eventType, Payload: jsonBytes}, nil)
	return err
}

func (m *consulEventDispatcher) RegisterEventHandler(handler interface{}) error {

	handlerType := reflect.TypeOf(handler)
	var eventType reflect.Type

	if handlerType.Kind() != reflect.Func || handlerType.NumIn() != 1 || handlerType.NumOut() != 1 {
		return fmt.Errorf("handler must be a func with one argument and return error")
	}

	if handlerType.NumIn() == 1 {
		eventType = handlerType.In(0)
	}

	waiter, err := newConsulEventWaiter(m.client, eventType.String(), eventType, handler)
	if err != nil {
		return err
	}

	m.waiters = append(m.waiters, waiter)
	waiter.WaitForEvents()

	return nil
}

type consulEventWaiter struct {
	client     *api.Client
	eventName  string
	eventType  reflect.Type
	handler    interface{}
	waitFunc   func()
	eventIDMap map[string]bool
}

func newConsulEventWaiter(client *api.Client, eventName string, eventType reflect.Type, handler interface{}) (*consulEventWaiter, error) {
	waiter := &consulEventWaiter{
		client:     client,
		eventName:  eventName,
		eventType:  eventType,
		handler:    handler,
		eventIDMap: map[string]bool{},
	}

	events, _, err := client.Event().List(eventName, &api.QueryOptions{AllowStale: true})
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		waiter.eventIDMap[event.ID] = true
	}

	return waiter, nil
}

func (c *consulEventWaiter) unserialiseEvent(serialised []byte) (reflect.Value, error) {
	event := reflect.New(c.eventType)
	err := json.Unmarshal(serialised, event.Interface())
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	return event.Elem(), nil
}

func consulWaiterRunFunc(method reflect.Value, event reflect.Value) (err error) {
	var result []reflect.Value

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("encountered a panic: %s", e)
		}
	}()

	result = method.Call([]reflect.Value{event})

	if result[0].Interface() != nil {
		var ok bool
		err, ok = result[0].Interface().(error)
		if !ok {
			err = fmt.Errorf("handler for '%s' returned non-error interface as first return value", reflect.TypeOf(event))
		}
	}
	return
}

func (c *consulEventWaiter) WaitForEvents() {

	c.waitFunc = func() {
		var lastIndex uint64

		for {
			events, queryMeta, err := c.client.Event().List(c.eventName, &api.QueryOptions{AllowStale: true, WaitIndex: lastIndex})
			if err != nil {
				fmt.Printf("error reading events from consul, error: %v", err)
				return
			}

			idMap := map[string]bool{}
			for _, e := range events {
				idMap[e.ID] = true
				if _, ok := c.eventIDMap[e.ID]; ok {
					continue
				}

				event, err := c.unserialiseEvent(e.Payload)
				if err != nil {
					log.Errorf("failed to unserialise event for '%s': %s", c.eventType, string(e.Payload))
					continue
				}

				method := reflect.ValueOf(c.handler)
				if err := consulWaiterRunFunc(method, event); err != nil {
					log.
						WithField("handler", runtime.FuncForPC(method.Pointer()).Name()).
						Errorf("handler for '%s' failed: %s", c.eventType, err)
				}
			}

			c.eventIDMap = idMap
			if queryMeta.LastIndex >= lastIndex {
				lastIndex = queryMeta.LastIndex
			} else if queryMeta.LastIndex == 0 {
				/*
					https://www.consul.io/api/features/blocking.html

					After the initial request (or a reset as above) the X-Consul-Index returned should always
					be greater than zero. It is a bug in Consul if it is not, however this has happened
					a few times and can still be triggered on some older Consul versions. It's especially bad because
					it causes blocking clients that are not aware to enter a busy loop, using excessive client CPU and
					causing high load on servers. It is always safe to use an index of 1 to wait for updates when the
					data being requested doesn't exist yet, so clients should sanity check that their index is at least
					1 after each blocking response is handled to be sure they actually block on the next request.
				*/
				lastIndex = 1
			} else {
				/*
					https://www.consul.io/api/features/blocking.html

					While indexes in general are monotonically increasing(i.e. they should only ever increase as time passes),
					there are several real-world scenarios in which they can go backwards for a given query.
					Implementations must check to see if a returned index is lower than the previous value, and if it is,
					should reset index to 0 - effectively restarting their blocking loop. Failure to do so may cause the client
					to miss future updates for an unbounded time, or to use an invalid index value that causes no blocking
					and increases load on the servers.
				*/
				lastIndex = 0
			}
		}
	}

	go c.waitFunc()

}

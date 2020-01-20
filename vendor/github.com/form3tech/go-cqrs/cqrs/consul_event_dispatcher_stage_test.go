package cqrs

import (
	"context"
	"fmt"
	"testing"

	"github.com/form3tech/go-cqrs/cqrs/support"
	log "github.com/sirupsen/logrus"

	"github.com/giantswarm/retry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type consulEventDispatcherStage struct {
	t                  *testing.T
	eventDispatcher    EventDispatcher
	result1            int
	result2            int
	registrationError  error
	registrationError2 error
	dispatchError      error
	logHook            *support.LogHook
}

type ConsulEvent struct {
	Number int
}

func ConsulEventDispatcherTest(t *testing.T) (*consulEventDispatcherStage, *consulEventDispatcherStage, *consulEventDispatcherStage) {

	d, err := newConsulEventDispatcher()
	require.NoError(t, err)
	s := &consulEventDispatcherStage{
		t:               t,
		eventDispatcher: d,
		logHook:         support.NewLogHook(),
	}

	log.AddHook(s.logHook)
	return s, s, s
}

func (s *consulEventDispatcherStage) an_event_handler_has_been_registered() *consulEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(func(e ConsulEvent) error {
		log.Printf("Handler ran!")
		s.result1 = e.Number
		return nil
	})

	return s
}

func (s *consulEventDispatcherStage) an_event_handler_has_been_registered_which_errors_with(message string) *consulEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(func(e ConsulEvent) error {
		s.result1 = e.Number
		return fmt.Errorf(message)
	})

	return s
}

func (s *consulEventDispatcherStage) an_event_handler_has_been_registered_which_panics_with(message string) *consulEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(func(e ConsulEvent) error {
		s.result1 = e.Number
		panic(message)
	})

	return s
}

func (s *consulEventDispatcherStage) a_second_event_handler_has_been_registered() *consulEventDispatcherStage {

	s.registrationError2 = s.eventDispatcher.RegisterEventHandler(func(e ConsulEvent) error {
		s.result2 = e.Number
		return nil
	})

	return s
}

func (s *consulEventDispatcherStage) the_event_is_dispatched() *consulEventDispatcherStage {

	ctx := context.Background()
	s.dispatchError = s.eventDispatcher.Dispatch(&ctx, ConsulEvent{Number: 17})

	return s
}

func (s *consulEventDispatcherStage) a_handler_that_is_not_a_function_is_registered() *consulEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(ConsulEvent{})

	return s
}

func (s *consulEventDispatcherStage) a_handler_with_more_than_one_argument_is_registered() *consulEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(func(a, b ConsulEvent) error { return nil })

	return s
}

func (s *consulEventDispatcherStage) a_handler_with_no_arguments_is_registered() *consulEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(func() error { return nil })

	return s
}

func (s *consulEventDispatcherStage) a_handler_with_more_than_one_return_value_is_registered() *consulEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(func(a ConsulEvent) (error, error) { return nil, nil })

	return s
}

func (s *consulEventDispatcherStage) a_handler_with_no_return_value_is_registered() *consulEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(func(a ConsulEvent) {})

	return s
}

func (s *consulEventDispatcherStage) the_event_handler_should_be_executed() *consulEventDispatcherStage {

	expected := 17
	err := retry.Do(func() error {
		if s.result1 != expected {
			return fmt.Errorf("expected result1 to be :%d but was: %d", expected, s.result1)
		}
		return nil
	}, stageRetryOptions...)

	assert.Nil(s.t, err)

	return s
}

func (s *consulEventDispatcherStage) the_second_event_handler_should_be_executed() *consulEventDispatcherStage {

	expected := 17
	err := retry.Do(func() error {
		if s.result2 != expected {
			return fmt.Errorf("expected result2 to be :%d but was: %d", expected, s.result1)
		}
		return nil
	}, stageRetryOptions...)

	assert.Nil(s.t, err)

	return s
}

func (s *consulEventDispatcherStage) and() *consulEventDispatcherStage {
	return s
}

func (s *consulEventDispatcherStage) no_error_should_be_returned() *consulEventDispatcherStage {
	assert.NoError(s.t, s.dispatchError)
	return s
}

func (s *consulEventDispatcherStage) a_registration_error_should_be_returned_with_message(errorMessage string) *consulEventDispatcherStage {
	require.Error(s.t, s.registrationError)
	assert.Equal(s.t, errorMessage, s.registrationError.Error())
	return s
}

func (s *consulEventDispatcherStage) an_error_should_be_logged_with_message(errorMessage string) *consulEventDispatcherStage {
	if err := retry.Do(func() error {
		for _, logLine := range s.logHook.Get() {
			if logLine.Message == errorMessage {
				return nil
			}
		}
		return fmt.Errorf("log line not found")
	}, stageRetryOptions...); err != nil {
		s.t.Fatalf("log line not found '%s': %s", errorMessage, err)
	}
	return s
}

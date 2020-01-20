package cqrs

import (
	"testing"

	"context"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type inMemoryEventDispatcherStage struct {
	t                  *testing.T
	eventDispatcher    EventDispatcher
	result1            int
	result2            int
	registrationError  error
	registrationError2 error
	dispatchError      error
	annotation         string
}

type testEvent struct {
	number int
}

func InMemoryEventDispatcherTest(t *testing.T) (*inMemoryEventDispatcherStage, *inMemoryEventDispatcherStage, *inMemoryEventDispatcherStage) {

	stage := &inMemoryEventDispatcherStage{
		t: t,
		eventDispatcher: &inMemoryEventDispatcher{
			handlers: make(map[string][]eventHandler),
		},
	}

	return stage, stage, stage
}

func (s *inMemoryEventDispatcherStage) an_event_handler_with_a_context_has_been_registered() *inMemoryEventDispatcherStage {
	s.registrationError = s.eventDispatcher.RegisterEventHandler(func(ctx *context.Context, e testEvent) error {
		s.result1 = e.number
		if ctx != nil {
			s.annotation = (*ctx).Value("handler").(string)
		}
		return nil
	})

	return s
}

func (s *inMemoryEventDispatcherStage) an_event_handler_has_been_registered() *inMemoryEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(func(e testEvent) error {
		s.result1 = e.number
		return nil
	})

	return s
}

func (s *inMemoryEventDispatcherStage) a_second_event_handler_has_been_registered() *inMemoryEventDispatcherStage {

	s.registrationError2 = s.eventDispatcher.RegisterEventHandler(func(e testEvent) error {
		s.result2 = e.number
		return nil
	})

	return s
}

func (s *inMemoryEventDispatcherStage) the_event_is_dispatched() *inMemoryEventDispatcherStage {
	ctx := context.Background()
	s.dispatchError = s.eventDispatcher.Dispatch(&ctx, testEvent{number: 11})

	return s

}

func (s *inMemoryEventDispatcherStage) the_event_handler_should_be_executed() *inMemoryEventDispatcherStage {

	assert.Equal(s.t, 11, s.result1)

	return s

}

func (s *inMemoryEventDispatcherStage) and() *inMemoryEventDispatcherStage {

	return s

}

func (s *inMemoryEventDispatcherStage) no_error_should_be_returned() *inMemoryEventDispatcherStage {

	assert.Nil(s.t, s.dispatchError)

	return s

}

func (s *inMemoryEventDispatcherStage) the_second_event_handler_should_be_executed() *inMemoryEventDispatcherStage {

	assert.Equal(s.t, 11, s.result2)

	return s
}

func (s *inMemoryEventDispatcherStage) a_handler_that_is_not_a_function_is_registered() *inMemoryEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(testEvent{})

	return s
}

func (s *inMemoryEventDispatcherStage) a_handler_that_with_more_than_one_argument_is_registered_without_using_context() *inMemoryEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(func(a, b testEvent) error { return nil })

	return s
}

func (s *inMemoryEventDispatcherStage) a_handler_that_with_no_arguments_is_registered() *inMemoryEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(func() error { return nil })

	return s
}

func (s *inMemoryEventDispatcherStage) a_handler_that_with_more_than_one_return_value_is_registered() *inMemoryEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(func(e testEvent) (int, error) { return 0, nil })

	return s
}

func (s *inMemoryEventDispatcherStage) a_handler_that_with_no_return_value_is_registered() *inMemoryEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(func(e testEvent) {})

	return s
}

func (s *inMemoryEventDispatcherStage) an_error_should_be_returned_with_message(errorMessage string) *inMemoryEventDispatcherStage {

	require.NotNil(s.t, s.dispatchError)
	assert.Error(s.t, s.dispatchError)
	assert.Equal(s.t, errorMessage, s.dispatchError.Error())

	return s
}

func (s *inMemoryEventDispatcherStage) a_registration_error_should_be_returned_with_message(errorMessage string) *inMemoryEventDispatcherStage {

	require.NotNil(s.t, s.registrationError)
	assert.Error(s.t, s.registrationError)
	assert.Equal(s.t, errorMessage, s.registrationError.Error())

	return s
}

func (s *inMemoryEventDispatcherStage) an_event_handler_has_been_registered_that_will_return_an_error(errorMessage string) *inMemoryEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(func(e testEvent) error {
		s.result1 = e.number
		return errors.New(errorMessage)
	})

	return s
}

func (s *inMemoryEventDispatcherStage) an_event_handler_has_been_registered_that_will_panic(panicMessage string) *inMemoryEventDispatcherStage {

	s.registrationError = s.eventDispatcher.RegisterEventHandler(func(e testEvent) error {
		s.result1 = e.number
		panic(panicMessage)
	})

	return s
}

func (s *inMemoryEventDispatcherStage) the_context_should_be_annotated_with_the_handler() *inMemoryEventDispatcherStage {
	assert.Equal(s.t, "github.com/form3tech/go-cqrs/cqrs.(*inMemoryEventDispatcherStage).an_event_handler_with_a_context_has_been_registered.func1", s.annotation)
	return s
}

package cqrs

import (
	"context"
	"fmt"
	"testing"

	"github.com/form3tech/go-cqrs/cqrs/support"
	"github.com/form3tech/go-messaging/messaging"
	"github.com/giantswarm/retry-go"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type queuedEventDispatcherStage struct {
	t                        *testing.T
	queuedEventDispatcher    *QueuedEventDispatcher
	commandHandlerCallCount  int
	commandHandler2CallCount int
	receivedEventValue       string
	receivedEvent2Value      string
	expectedEventValue       string
	registrationError        error
	registrationError2       error
	executeError             error
	receiver                 messaging.Receiver
	logHook                  *support.LogHook
	expectedError            string
	annotation               string
}

type testQueuedEvent struct{ Value string }

func QueuedEventDispatcherTest(t *testing.T) (*queuedEventDispatcherStage, *queuedEventDispatcherStage, *queuedEventDispatcherStage) {

	sender := messaging.NewSqsSender()
	receiver := messaging.NewSqsReceiverBuilder().
		WithVisibilityTimeout(3).
		Build()

	stage := &queuedEventDispatcherStage{
		t:                     t,
		queuedEventDispatcher: newQueuedEventDispatcher("test", receiver, sender),
		receiver:              receiver,
		logHook:               support.NewLogHook(),
		expectedError:         uuid.New().String(),
	}

	log.AddHook(stage.logHook)

	log.SetLevel(log.DebugLevel)

	support.PurgeAllSQS()
	stage.expectedEventValue = uuid.New().String()

	return stage, stage, stage
}

func (s *queuedEventDispatcherStage) TearDown() *queuedEventDispatcherStage {
	s.receiver.Close()
	return s
}

func (s *queuedEventDispatcherStage) the_event_is_queued() *queuedEventDispatcherStage {
	s.executeError = s.queuedEventDispatcher.Dispatch(nil, testQueuedEvent{Value: s.expectedEventValue})
	return s
}

func (s *queuedEventDispatcherStage) the_event_handler_should_be_executed_once() *queuedEventDispatcherStage {
	if err := retry.Do(func() error {
		if s.commandHandlerCallCount != 1 {
			return fmt.Errorf("not called once")
		}
		return nil
	}, stageRetryOptions...); err != nil {
		s.t.Fatalf("command was not run once - actually run %d time(s)", s.commandHandlerCallCount)
	}
	return s
}

func (s *queuedEventDispatcherStage) the_second_event_handler_should_be_executed_once() *queuedEventDispatcherStage {
	if err := retry.Do(func() error {
		if s.commandHandler2CallCount != 1 {
			return fmt.Errorf("not called once")
		}
		return nil
	}, stageRetryOptions...); err != nil {
		s.t.Fatalf("command was not run once - actually run %d time(s)", s.commandHandler2CallCount)
	}
	return s
}

func (s *queuedEventDispatcherStage) the_event_handler_should_have_received_the_correct_event_data() *queuedEventDispatcherStage {
	assert.Equal(s.t, s.expectedEventValue, s.receivedEventValue)
	return s
}

func (s *queuedEventDispatcherStage) the_second_event_handler_should_have_received_the_correct_event_data() *queuedEventDispatcherStage {
	assert.Equal(s.t, s.expectedEventValue, s.receivedEvent2Value)
	return s
}

func (s *queuedEventDispatcherStage) and() *queuedEventDispatcherStage {
	return s
}

func (s *queuedEventDispatcherStage) no_registration_error_should_be_returned() *queuedEventDispatcherStage {
	assert.NoError(s.t, s.registrationError)
	return s
}

func (s *queuedEventDispatcherStage) no_execution_error_should_be_returned() *queuedEventDispatcherStage {
	assert.NoError(s.t, s.executeError)
	return s
}

func (s *queuedEventDispatcherStage) the_error_should_be_logged() *queuedEventDispatcherStage {
	if err := retry.Do(func() error {
		for _, logLine := range s.logHook.Get() {
			if logLine.Message == fmt.Sprintf("error in handler #0 for 'cqrs.testQueuedEvent': %s", s.expectedError) {
				return nil
			}
		}
		return fmt.Errorf("log line not found")
	}, stageRetryOptions...); err != nil {
		s.t.Fatalf("log line not found '%s': %s", s.expectedError, err)
	}
	return s
}

func (s *queuedEventDispatcherStage) the_panic_should_be_logged() *queuedEventDispatcherStage {
	if err := retry.Do(func() error {
		for _, logLine := range s.logHook.Get() {
			if logLine.Message == fmt.Sprintf("error in handler #0 for 'cqrs.testQueuedEvent': encountered a panic: %s", s.expectedError) {
				return nil
			}
		}
		return fmt.Errorf("log line not found")
	}, stageRetryOptions...); err != nil {
		s.t.Fatalf("log line not found '%s': %s", s.expectedError, err)
	}
	return s
}

func (s *queuedEventDispatcherStage) a_second_event_handler_is_registered_for_the_same_command() *queuedEventDispatcherStage {
	s.registrationError2 = s.queuedEventDispatcher.RegisterEventHandler(func(c testQueuedEvent) error {
		s.commandHandler2CallCount++
		s.receivedEvent2Value = c.Value
		return nil
	})

	return s
}

func (s *queuedEventDispatcherStage) an_event_handler_is_registered_for_an_event() *queuedEventDispatcherStage {
	s.registrationError = s.queuedEventDispatcher.RegisterEventHandler(func(c testQueuedEvent) error {
		s.commandHandlerCallCount++
		s.receivedEventValue = c.Value
		return nil
	})
	return s
}

func (s *queuedEventDispatcherStage) an_event_handler_which_uses_a_context_is_registered_for_an_event() *queuedEventDispatcherStage {
	s.registrationError = s.queuedEventDispatcher.RegisterEventHandler(func(ctx *context.Context, c testQueuedEvent) error {
		s.commandHandlerCallCount++
		s.receivedEventValue = c.Value
		if ctx != nil {
			s.annotation = (*ctx).Value("handler").(string)
		}
		return nil
	})
	return s
}

func (s *queuedEventDispatcherStage) an_event_handler_which_errors_is_registered_for_a_command() *queuedEventDispatcherStage {
	s.registrationError = s.queuedEventDispatcher.RegisterEventHandler(func(c testQueuedEvent) error {
		s.commandHandlerCallCount++
		s.receivedEventValue = c.Value

		return fmt.Errorf(s.expectedError)
	})
	return s
}

func (s *queuedEventDispatcherStage) an_event_handler_which_panics_is_registered_for_a_command() *queuedEventDispatcherStage {
	s.registrationError = s.queuedEventDispatcher.RegisterEventHandler(func(c testQueuedEvent) error {
		s.commandHandlerCallCount++
		s.receivedEventValue = c.Value
		panic(s.expectedError)
	})
	return s
}

func (s *queuedEventDispatcherStage) the_event_handler_should_be_executed() *queuedEventDispatcherStage {
	return s.no_registration_error_should_be_returned().and().
		no_execution_error_should_be_returned().and().
		the_event_handler_should_be_executed_once().and().
		the_event_handler_should_have_received_the_correct_event_data()
}

func (s *queuedEventDispatcherStage) the_second_event_handler_should_be_executed() *queuedEventDispatcherStage {
	return s.the_second_event_handler_should_be_executed_once().and().
		the_second_event_handler_should_have_received_the_correct_event_data()
}

func (s *queuedEventDispatcherStage) the_context_should_be_annotated_with_the_handler() *queuedEventDispatcherStage {
	assert.Equal(s.t, "github.com/form3tech/go-cqrs/cqrs.(*queuedEventDispatcherStage).an_event_handler_which_uses_a_context_is_registered_for_an_event.func1", s.annotation)
	return s
}

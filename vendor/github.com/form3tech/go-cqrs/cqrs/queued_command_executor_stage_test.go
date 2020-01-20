package cqrs

import (
	"context"
	"fmt"
	"testing"

	"github.com/form3tech/go-cqrs/cqrs/support"
	"github.com/form3tech/go-messaging/messaging"
	"github.com/giantswarm/retry-go"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type queuedCommandExecutorStage struct {
	t                        *testing.T
	queuedCommandExecutor    *QueuedCommandExecutor
	commandHandlerCallCount  int
	receivedCommandValue     string
	commandHandlerCallCount2 int
	receivedCommandValue2    string
	expectedCommandValue     string
	registrationError        error
	registrationError2       error
	executeError             error
	executeError2            error
	receiver                 messaging.Receiver
	logHook                  *support.LogHook
	expectedError            string
	annotation               string
}

type testQueuedCommand struct{ Value string }
type testQueuedCommand2 struct{ Value string }

func QueuedCommandExecutorTest(t *testing.T) (*queuedCommandExecutorStage, *queuedCommandExecutorStage, *queuedCommandExecutorStage) {

	sender := messaging.NewSqsSender()
	receiver := messaging.NewSqsReceiverBuilder().
		WithVisibilityTimeout(3).
		Build()

	stage := &queuedCommandExecutorStage{
		t:                     t,
		queuedCommandExecutor: newQueuedCommandExecutor(&sqlx.DB{}, "test", receiver, sender),
		receiver:              receiver,
		logHook:               support.NewLogHook(),
		expectedError:         uuid.New().String(),
	}

	log.AddHook(stage.logHook)

	log.SetLevel(log.DebugLevel)

	support.PurgeAllSQS()
	stage.expectedCommandValue = uuid.New().String()

	return stage, stage, stage
}

func (s *queuedCommandExecutorStage) TearDown() *queuedCommandExecutorStage {
	s.receiver.Close()
	return s
}

func (s *queuedCommandExecutorStage) the_command_is_executed() *queuedCommandExecutorStage {
	id := uuid.New()
	s.executeError = s.queuedCommandExecutor.Execute(nil, &id, testQueuedCommand{Value: s.expectedCommandValue})
	return s
}

func (s *queuedCommandExecutorStage) both_commands_are_executed() *queuedCommandExecutorStage {
	id := uuid.New()
	s.executeError = s.queuedCommandExecutor.Execute(nil, &id, testQueuedCommand{Value: s.expectedCommandValue})
	s.executeError2 = s.queuedCommandExecutor.Execute(nil, &id, testQueuedCommand2{Value: s.expectedCommandValue})
	return s
}

func (s *queuedCommandExecutorStage) the_command_is_executed_with_context() *queuedCommandExecutorStage {
	id := uuid.New()
	ctx := context.Background()
	s.executeError = s.queuedCommandExecutor.Execute(&ctx, &id, testQueuedCommand{Value: s.expectedCommandValue})
	return s
}

func (s *queuedCommandExecutorStage) the_command_handler_should_be_executed_once() *queuedCommandExecutorStage {
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

func (s *queuedCommandExecutorStage) the_command_handler_should_have_received_the_correct_command_data() *queuedCommandExecutorStage {
	assert.Equal(s.t, s.expectedCommandValue, s.receivedCommandValue)
	return s
}

func (s *queuedCommandExecutorStage) and() *queuedCommandExecutorStage {
	return s
}

func (s *queuedCommandExecutorStage) no_registration_error_should_be_returned() *queuedCommandExecutorStage {
	assert.NoError(s.t, s.registrationError)
	return s
}

func (s *queuedCommandExecutorStage) no_registration_error_should_be_returned_for_either_command() *queuedCommandExecutorStage {
	assert.NoError(s.t, s.registrationError)
	assert.NoError(s.t, s.registrationError2)
	return s
}

func (s *queuedCommandExecutorStage) no_execution_error_should_be_returned() *queuedCommandExecutorStage {
	assert.NoError(s.t, s.executeError)
	return s
}

func (s *queuedCommandExecutorStage) an_execution_error_should_be_returned_with_message(message string) *queuedCommandExecutorStage {
	assert.EqualError(s.t, s.executeError, message)
	return s
}

func (s *queuedCommandExecutorStage) the_error_should_be_logged() *queuedCommandExecutorStage {
	if err := retry.Do(func() error {
		for _, logLine := range s.logHook.Get() {
			if logLine.Message == s.expectedError {
				return nil
			}
		}
		return fmt.Errorf("log line not found")
	}, stageRetryOptions...); err != nil {
		s.t.Fatalf("log line not found '%s': %s", s.expectedError, err)
	}
	return s
}

func (s *queuedCommandExecutorStage) the_panic_should_be_logged() *queuedCommandExecutorStage {
	if err := retry.Do(func() error {
		for _, logLine := range s.logHook.Get() {
			if logLine.Message == fmt.Sprintf("handler for 'cqrs.testQueuedCommand' encountered a panic: %s", s.expectedError) {
				return nil
			}
		}
		return fmt.Errorf("log line not found")
	}, stageRetryOptions...); err != nil {
		s.t.Fatalf("log line not found '%s': %s", s.expectedError, err)
	}
	return s
}

func (s *queuedCommandExecutorStage) another_command_handler_is_registered_for_the_same_command() *queuedCommandExecutorStage {
	s.registrationError2 = s.queuedCommandExecutor.RegisterCommandHandler(func(c testQueuedCommand) error {
		return nil
	}, AllowEveryone())

	return s
}

func (s *queuedCommandExecutorStage) another_command_handler_is_registered_for_another_command() *queuedCommandExecutorStage {
	s.registrationError2 = s.queuedCommandExecutor.RegisterCommandHandler(func(c testQueuedCommand2) error {
		s.commandHandlerCallCount2++
		s.receivedCommandValue2 = c.Value
		return nil
	}, AllowEveryone())

	return s
}

func (s *queuedCommandExecutorStage) the_second_command_returns_an_error_on_registration_with_message(message string) *queuedCommandExecutorStage {
	assert.EqualError(s.t, s.registrationError2, message)
	return s
}

func (s *queuedCommandExecutorStage) the_executor_uses_a_shared_naming_strategy() *queuedCommandExecutorStage {
	s.queuedCommandExecutor.QueueNamer = SharedQueueNamingStrategy{ApplicationName: "test", QueueName: "commands"}
	return s
}

func (s *queuedCommandExecutorStage) a_command_handler_is_registered_for_a_command() *queuedCommandExecutorStage {
	s.registrationError = s.queuedCommandExecutor.RegisterCommandHandler(func(c testQueuedCommand) error {
		s.commandHandlerCallCount++
		s.receivedCommandValue = c.Value

		return nil
	}, AllowEveryone())
	return s
}

func (s *queuedCommandExecutorStage) a_command_handler_which_errors_is_registered_for_a_command() *queuedCommandExecutorStage {
	s.registrationError = s.queuedCommandExecutor.RegisterCommandHandler(func(c testQueuedCommand) error {
		s.commandHandlerCallCount++
		s.receivedCommandValue = c.Value

		return fmt.Errorf(s.expectedError)
	}, AllowEveryone())
	return s
}

func (s *queuedCommandExecutorStage) a_command_handler_which_panics_is_registered_for_a_command() *queuedCommandExecutorStage {
	s.registrationError = s.queuedCommandExecutor.RegisterCommandHandler(func(c testQueuedCommand) error {
		s.commandHandlerCallCount++
		s.receivedCommandValue = c.Value
		panic(s.expectedError)
	}, AllowEveryone())
	return s
}

func (s *queuedCommandExecutorStage) a_command_handler_that_uses_a_database_connection_is_registered() *queuedCommandExecutorStage {
	s.registrationError = s.queuedCommandExecutor.RegisterCommandHandler(func(ctx *context.Context, db *sqlx.DB, c testQueuedCommand) error {
		if db == nil {
			return errors.New("db is nil in command handler")
		}

		s.commandHandlerCallCount++
		s.receivedCommandValue = c.Value

		return nil
	}, AllowEveryone())

	return s
}

func (s *queuedCommandExecutorStage) a_command_handler_that_uses_a_context_and_database_is_registered() *queuedCommandExecutorStage {
	s.registrationError = s.queuedCommandExecutor.RegisterCommandHandler(func(ctx *context.Context, db *sqlx.DB, c testQueuedCommand) error {
		if ctx == nil {
			return errors.New("context is nil in command handler")
		}

		if db == nil {
			return errors.New("db is nil in command handler")
		}

		s.commandHandlerCallCount++
		s.receivedCommandValue = c.Value

		return nil
	}, AllowEveryone())

	return s
}

func (s *queuedCommandExecutorStage) a_command_handler_that_uses_a_context_is_registered() *queuedCommandExecutorStage {
	s.registrationError = s.queuedCommandExecutor.RegisterCommandHandler(func(ctx *context.Context, db *sqlx.DB, c testQueuedCommand) error {
		if ctx == nil {
			return errors.New("context is nil in command handler")
		}

		s.annotation = (*ctx).Value("handler").(string)

		s.commandHandlerCallCount++
		s.receivedCommandValue = c.Value

		return nil
	}, AllowEveryone())

	return s
}

func (s *queuedCommandExecutorStage) the_command_handler_should_be_executed() *queuedCommandExecutorStage {
	return s.no_registration_error_should_be_returned().and().
		no_execution_error_should_be_returned().and().
		the_command_handler_should_be_executed_once().and().
		the_command_handler_should_have_received_the_correct_command_data()
}

func (s *queuedCommandExecutorStage) both_command_handlers_should_be_executed() *queuedCommandExecutorStage {
	s.the_command_handler_should_be_executed()

	assert.NoError(s.t, s.executeError2)
	if err := retry.Do(func() error {
		if s.commandHandlerCallCount2 != 1 {
			return fmt.Errorf("not called once")
		}
		return nil
	}, stageRetryOptions...); err != nil {
		s.t.Fatalf("command was not run once - actually run %d time(s)", s.commandHandlerCallCount)
	}
	assert.Equal(s.t, s.expectedCommandValue, s.receivedCommandValue2)
	return s
}

func (s *queuedCommandExecutorStage) the_context_should_be_annotated_with_the_handler() *queuedCommandExecutorStage {
	assert.Equal(s.t, "github.com/form3tech/go-cqrs/cqrs.(*queuedCommandExecutorStage).a_command_handler_that_uses_a_context_is_registered.func1", s.annotation)
	return s
}

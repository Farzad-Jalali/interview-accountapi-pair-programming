package cqrs

import (
	"fmt"
	"testing"

	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type commandExecutorStage struct {
	t                    *testing.T
	registrationError    error
	executionError       error
	commandExecutor      CommandExecutor
	commandHandlerCalled bool
	commandValue         int
	annotation           string
}

type testCommand struct {
	value int
}

func AllowEveryone() func(*context.Context, *uuid.UUID) error {
	return func(ctx *context.Context, organisationId *uuid.UUID) error {
		return nil
	}
}

func CommandExecutorTest(t *testing.T) (*commandExecutorStage, *commandExecutorStage, *commandExecutorStage) {

	stage := &commandExecutorStage{
		t: t,
		commandExecutor: &inMemoryCommandExecutor{
			handlers: make(map[string]h),
			db:       &sqlx.DB{},
		},
		commandHandlerCalled: false,
	}

	return stage, stage, stage
}

func (s *commandExecutorStage) a_command_handler_is_registered() *commandExecutorStage {

	s.registrationError = s.commandExecutor.RegisterCommandHandler(func(c testCommand) error {

		s.commandHandlerCalled = true
		s.commandValue = c.value

		return nil

	}, AllowEveryone())

	return s
}

func (s *commandExecutorStage) a_command_handler_which_returns_an_error_is_registered_for_test_command() *commandExecutorStage {

	s.registrationError = s.commandExecutor.RegisterCommandHandler(func(c testCommand) error {

		s.commandHandlerCalled = true
		s.commandValue = c.value

		return fmt.Errorf("oh no: a test error")

	}, AllowEveryone())

	return s
}

func (s *commandExecutorStage) a_command_handler_which_panics_is_registered_for_test_command() *commandExecutorStage {
	s.registrationError = s.commandExecutor.RegisterCommandHandler(func(c testCommand) error {
		s.commandHandlerCalled = true
		s.commandValue = c.value
		panic("oh no: a test panic occurred")
	}, AllowEveryone())

	return s
}

func (s *commandExecutorStage) a_command_handler_that_uses_a_database_connection_is_registered() *commandExecutorStage {
	// Should i add context here ??
	s.registrationError = s.commandExecutor.RegisterCommandHandler(func(ctx *context.Context, db *sqlx.DB, c testCommand) error {
		if db != nil {
			s.commandHandlerCalled = true
			s.commandValue = c.value
		}

		return nil
	}, AllowEveryone())

	return s
}

func (s *commandExecutorStage) another_command_handler_is_registered_for_the_same_command() *commandExecutorStage {
	s.registrationError = s.commandExecutor.RegisterCommandHandler(func(c testCommand) error {
		s.commandHandlerCalled = true
		s.commandValue = c.value

		return nil
	}, AllowEveryone())

	return s
}

func (s *commandExecutorStage) the_command_is_executed() *commandExecutorStage {
	id := uuid.New()
	s.executionError = s.commandExecutor.Execute(nil, &id, testCommand{value: 42})

	return s
}

func (s *commandExecutorStage) the_command_is_executed_with_context() *commandExecutorStage {
	id := uuid.New()
	ctx := context.Background()
	s.executionError = s.commandExecutor.Execute(&ctx, &id, testCommand{value: 42})

	return s
}

func (s *commandExecutorStage) the_command_handler_should_be_executed() *commandExecutorStage {
	assert.True(s.t, s.commandHandlerCalled)
	assert.Equal(s.t, 42, s.commandValue)

	return s
}

func (s *commandExecutorStage) and() *commandExecutorStage {

	return s
}

func (s *commandExecutorStage) no_error_should_be_returned() *commandExecutorStage {

	assert.Nil(s.t, s.registrationError)
	assert.Nil(s.t, s.executionError)

	return s
}

func (s *commandExecutorStage) a_registration_error_is_returned_with_message(message string) *commandExecutorStage {
	assert.EqualError(s.t, s.registrationError, message)

	return s
}

func (s *commandExecutorStage) an_execution_error_is_returned_with_message(message string) *commandExecutorStage {
	assert.EqualError(s.t, s.executionError, message)

	return s
}

func (s *commandExecutorStage) a_command_handler_that_uses_a_context_and_database_is_registered() *commandExecutorStage {

	s.registrationError = s.commandExecutor.RegisterCommandHandler(func(ctx *context.Context, db *sqlx.DB, c testCommand) error {

		if ctx != nil && db != nil {
			s.commandHandlerCalled = true
			s.commandValue = c.value
		}

		return nil

	}, AllowEveryone())

	return s
}

func (s *commandExecutorStage) a_command_handler_that_uses_a_context_is_registered() *commandExecutorStage {

	s.registrationError = s.commandExecutor.RegisterCommandHandler(func(ctx *context.Context, db *sqlx.DB, c testCommand) error {

		if ctx != nil {
			s.commandHandlerCalled = true
			s.commandValue = c.value
			s.annotation = (*ctx).Value("handler").(string)
		}

		return nil

	}, AllowEveryone())

	return s
}

func (s *commandExecutorStage) the_context_should_be_annotated_with_the_handler() *commandExecutorStage {
	assert.Equal(s.t, "github.com/form3tech/go-cqrs/cqrs.(*commandExecutorStage).a_command_handler_that_uses_a_context_is_registered.func1", s.annotation)
	return s
}

package cqrs

import (
	"testing"
)

func Test_execute_command_handler(t *testing.T) {

	given, when, then := CommandExecutorTest(t)

	given.
		a_command_handler_is_registered()

	when.
		the_command_is_executed()

	then.
		the_command_handler_should_be_executed().and().
		no_error_should_be_returned()

}

func Test_registering_two_handlers_for_same_command(t *testing.T) {

	given, when, then := CommandExecutorTest(t)

	given.
		a_command_handler_is_registered()

	when.
		another_command_handler_is_registered_for_the_same_command()

	then.
		a_registration_error_is_returned_with_message("handler already registered for cqrs.testCommand only one handler allowed per command")

}

func Test_executing_unregistered_command(t *testing.T) {

	_, when, then := CommandExecutorTest(t)

	when.
		the_command_is_executed()

	then.
		an_execution_error_is_returned_with_message("no command handler registered of type: cqrs.testCommand")

}

func Test_registered_command_handler_that_uses_database(t *testing.T) {

	given, when, then := CommandExecutorTest(t)

	given.
		a_command_handler_that_uses_a_database_connection_is_registered()

	when.
		the_command_is_executed_with_context()

	then.
		the_command_handler_should_be_executed().and().
		no_error_should_be_returned()
}

func Test_registered_command_handler_that_uses_context(t *testing.T) {

	given, when, then := CommandExecutorTest(t)

	given.
		a_command_handler_that_uses_a_context_is_registered()

	when.
		the_command_is_executed_with_context()

	then.
		the_command_handler_should_be_executed().and().
		no_error_should_be_returned().and().
		the_context_should_be_annotated_with_the_handler()
}

func Test_registered_command_handler_that_uses_context_and_database(t *testing.T) {

	given, when, then := CommandExecutorTest(t)

	given.
		a_command_handler_that_uses_a_context_and_database_is_registered()

	when.
		the_command_is_executed_with_context()

	then.
		the_command_handler_should_be_executed().and().
		no_error_should_be_returned()
}

func Test_execute_command_handler_which_returns_error(t *testing.T) {

	given, when, then := CommandExecutorTest(t)

	given.
		a_command_handler_which_returns_an_error_is_registered_for_test_command()

	when.
		the_command_is_executed_with_context()

	then.
		the_command_handler_should_be_executed().and().
		an_execution_error_is_returned_with_message("oh no: a test error")

}

func Test_execute_command_handler_which_panics(t *testing.T) {

	given, when, then := CommandExecutorTest(t)

	given.
		a_command_handler_which_panics_is_registered_for_test_command()

	when.
		the_command_is_executed_with_context()

	then.
		the_command_handler_should_be_executed().and().
		an_execution_error_is_returned_with_message("handler for 'cqrs.testCommand' encountered a panic: oh no: a test panic occurred")

}

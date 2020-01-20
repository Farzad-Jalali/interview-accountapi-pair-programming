package cqrs

import "testing"

func Test_execute_queued_command_handler(t *testing.T) {

	given, when, then := QueuedCommandExecutorTest(t)
	defer then.TearDown()

	given.
		a_command_handler_is_registered_for_a_command()

	when.
		the_command_is_executed()

	then.
		no_registration_error_should_be_returned().and().
		no_execution_error_should_be_returned().and().
		the_command_handler_should_be_executed_once().and().
		the_command_handler_should_have_received_the_correct_command_data()

}

func Test_registered_queued_command_handler_that_uses_database(t *testing.T) {

	given, when, then := QueuedCommandExecutorTest(t)
	defer then.TearDown()

	given.
		a_command_handler_that_uses_a_database_connection_is_registered()

	when.
		the_command_is_executed_with_context()

	then.
		the_command_handler_should_be_executed()
}

func Test_registering_two_queued_handlers_for_same_command(t *testing.T) {

	given, when, then := QueuedCommandExecutorTest(t)
	defer then.TearDown()

	given.
		a_command_handler_is_registered_for_a_command()

	when.
		another_command_handler_is_registered_for_the_same_command()

	then.
		no_registration_error_should_be_returned().and().
		the_second_command_returns_an_error_on_registration_with_message("handler already registered for cqrs.testQueuedCommand only one handler allowed per command")
}

func Test_executing_command_without_a_queued_command_handler(t *testing.T) {

	_, when, then := QueuedCommandExecutorTest(t)
	defer then.TearDown()

	when.
		the_command_is_executed()

	then.
		an_execution_error_should_be_returned_with_message("no command handler registered of type: cqrs.testQueuedCommand")
}

func Test_registered_queued_command_handler_that_uses_context(t *testing.T) {

	given, when, then := QueuedCommandExecutorTest(t)
	defer then.TearDown()

	given.
		a_command_handler_that_uses_a_context_is_registered()

	when.
		the_command_is_executed_with_context()

	then.
		the_command_handler_should_be_executed().and().
		the_context_should_be_annotated_with_the_handler()
}

func Test_registered_queued_command_handler_that_uses_context_and_database(t *testing.T) {

	given, when, then := QueuedCommandExecutorTest(t)
	defer then.TearDown()

	given.
		a_command_handler_that_uses_a_context_and_database_is_registered()

	when.
		the_command_is_executed_with_context()

	then.
		the_command_handler_should_be_executed()
}

func Test_execute_queued_command_handler_which_returns_error(t *testing.T) {

	given, when, then := QueuedCommandExecutorTest(t)
	defer then.TearDown()

	given.
		a_command_handler_which_errors_is_registered_for_a_command()

	when.
		the_command_is_executed()

	then.
		the_command_handler_should_be_executed().and().
		the_error_should_be_logged()

}

func Test_execute_queued_command_handler_which_panics(t *testing.T) {

	given, when, then := QueuedCommandExecutorTest(t)
	defer then.TearDown()

	given.
		a_command_handler_which_panics_is_registered_for_a_command()

	when.
		the_command_is_executed()

	then.
		the_command_handler_should_be_executed().and().
		the_panic_should_be_logged()

}

func Test_register_queued_command_handlers_for_two_different_commands_with_shared_naming_strategy(t *testing.T) {

	given, when, then := QueuedCommandExecutorTest(t)
	defer then.TearDown()

	given.
		the_executor_uses_a_shared_naming_strategy().and().
		a_command_handler_is_registered_for_a_command()

	when.
		another_command_handler_is_registered_for_another_command()

	then.
		no_registration_error_should_be_returned_for_either_command()

}

func Test_execute_queued_command_handlers_for_two_different_commands_with_shared_naming_strategy(t *testing.T) {

	given, when, then := QueuedCommandExecutorTest(t)
	defer then.TearDown()

	given.
		the_executor_uses_a_shared_naming_strategy().and().
		a_command_handler_is_registered_for_a_command().and().
		another_command_handler_is_registered_for_another_command()

	when.
		both_commands_are_executed()

	then.
		both_command_handlers_should_be_executed()

}

package cqrs

import "testing"

func Test_dispatch_queued_event_handler(t *testing.T) {

	given, when, then := QueuedEventDispatcherTest(t)
	defer then.TearDown()

	given.
		an_event_handler_is_registered_for_an_event()

	when.
		the_event_is_queued()

	then.
		no_registration_error_should_be_returned().and().
		no_execution_error_should_be_returned().and().
		the_event_handler_should_be_executed_once().and().
		the_event_handler_should_have_received_the_correct_event_data()

}

func Test_dispatch_queued_event_handler_with_context(t *testing.T) {

	given, when, then := QueuedEventDispatcherTest(t)
	defer then.TearDown()

	given.
		an_event_handler_which_uses_a_context_is_registered_for_an_event()

	when.
		the_event_is_queued()

	then.
		no_registration_error_should_be_returned().and().
		no_execution_error_should_be_returned().and().
		the_event_handler_should_be_executed_once().and().
		the_event_handler_should_have_received_the_correct_event_data().and().
		the_context_should_be_annotated_with_the_handler()

}

func Test_registering_two_queued_handlers_for_same_event(t *testing.T) {

	given, when, then := QueuedEventDispatcherTest(t)
	defer then.TearDown()

	given.
		an_event_handler_is_registered_for_an_event().and().
		a_second_event_handler_is_registered_for_the_same_command()

	when.
		the_event_is_queued()

	then.
		the_event_handler_should_be_executed().and().
		the_second_event_handler_should_be_executed()

}

func Test_executing_command_without_a_queued_event_handler(t *testing.T) {

	_, when, then := QueuedEventDispatcherTest(t)
	defer then.TearDown()

	when.
		the_event_is_queued()

	then.
		no_execution_error_should_be_returned()
}

func Test_dispatch_queued_event_handler_which_returns_error(t *testing.T) {

	given, when, then := QueuedEventDispatcherTest(t)
	defer then.TearDown()

	given.
		an_event_handler_which_errors_is_registered_for_a_command()

	when.
		the_event_is_queued()

	then.
		the_event_handler_should_be_executed().and().
		the_error_should_be_logged()

}

func Test_dispatch_queued_event_handler_which_panics(t *testing.T) {

	given, when, then := QueuedEventDispatcherTest(t)
	defer then.TearDown()

	given.
		an_event_handler_which_panics_is_registered_for_a_command()

	when.
		the_event_is_queued()

	then.
		the_event_handler_should_be_executed().and().
		the_panic_should_be_logged()

}

func Test_dispatch_queued_event_handler_which_panics_and_another_which_does_not(t *testing.T) {

	given, when, then := QueuedEventDispatcherTest(t)
	defer then.TearDown()

	given.
		an_event_handler_which_panics_is_registered_for_a_command().and().
		a_second_event_handler_is_registered_for_the_same_command()

	when.
		the_event_is_queued()

	then.
		the_event_handler_should_be_executed().and().
		the_second_event_handler_should_be_executed().and().
		the_panic_should_be_logged()

}

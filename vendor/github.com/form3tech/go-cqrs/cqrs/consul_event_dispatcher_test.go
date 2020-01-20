package cqrs

import "testing"

func Test_ConsulEventDispatcher_dispatch_event_to_single_handler(t *testing.T) {

	given, when, then := ConsulEventDispatcherTest(t)

	given.
		an_event_handler_has_been_registered()

	when.
		the_event_is_dispatched()

	then.
		the_event_handler_should_be_executed().and().
		no_error_should_be_returned()

}

func Test_ConsulEventDispatcher_dispatch_event_to_multiple_handlers(t *testing.T) {

	given, when, then := ConsulEventDispatcherTest(t)

	given.
		an_event_handler_has_been_registered().and().
		a_second_event_handler_has_been_registered()

	when.
		the_event_is_dispatched()

	then.
		the_event_handler_should_be_executed().and().
		the_second_event_handler_should_be_executed().and().
		no_error_should_be_returned()
}

func Test_ConsulEventDispatcher_register_an_event_handler_that_is_not_a_function(t *testing.T) {

	_, when, then := ConsulEventDispatcherTest(t)

	when.
		a_handler_that_is_not_a_function_is_registered()

	then.
		a_registration_error_should_be_returned_with_message("handler must be a func with one argument and return error")

}

func Test_ConsulEventDispatcher_register_an_event_handler_with_more_than_one_argument(t *testing.T) {

	_, when, then := ConsulEventDispatcherTest(t)

	when.
		a_handler_with_more_than_one_argument_is_registered()

	then.
		a_registration_error_should_be_returned_with_message("handler must be a func with one argument and return error")

}

func Test_ConsulEventDispatcher_register_an_event_handler_with_no_arguments(t *testing.T) {

	_, when, then := ConsulEventDispatcherTest(t)

	when.
		a_handler_with_no_arguments_is_registered()

	then.
		a_registration_error_should_be_returned_with_message("handler must be a func with one argument and return error")

}

func Test_ConsulEventDispatcher_register_an_event_handler_with_more_than_one_return_value(t *testing.T) {

	_, when, then := ConsulEventDispatcherTest(t)

	when.
		a_handler_with_more_than_one_return_value_is_registered()

	then.
		a_registration_error_should_be_returned_with_message("handler must be a func with one argument and return error")

}

func Test_ConsulEventDispatcher_register_an_event_handler_with_no_return_value(t *testing.T) {

	_, when, then := ConsulEventDispatcherTest(t)

	when.
		a_handler_with_no_return_value_is_registered()

	then.
		a_registration_error_should_be_returned_with_message("handler must be a func with one argument and return error")

}

func Test_ConsulEventDispatcher_dispatch_event_to_single_handler_which_errors(t *testing.T) {

	given, when, then := ConsulEventDispatcherTest(t)

	given.
		an_event_handler_has_been_registered_which_errors_with("oh no")

	when.
		the_event_is_dispatched()

	then.
		the_event_handler_should_be_executed().and().
		an_error_should_be_logged_with_message("handler for 'cqrs.ConsulEvent' failed: oh no")

}

func Test_ConsulEventDispatcher_dispatch_event_to_single_handler_which_panics(t *testing.T) {

	given, when, then := ConsulEventDispatcherTest(t)

	given.
		an_event_handler_has_been_registered_which_panics_with("oh no")

	when.
		the_event_is_dispatched()

	then.
		the_event_handler_should_be_executed().and().
		an_error_should_be_logged_with_message("handler for 'cqrs.ConsulEvent' failed: encountered a panic: oh no")

}

func Test_ConsulEventDispatcher_dispatch_event_to_single_handler_which_errors_and_second_handler_which_does_not(t *testing.T) {

	given, when, then := ConsulEventDispatcherTest(t)

	given.
		an_event_handler_has_been_registered_which_errors_with("oh no").and().
		a_second_event_handler_has_been_registered()

	when.
		the_event_is_dispatched()

	then.
		the_event_handler_should_be_executed().and().
		the_second_event_handler_should_be_executed().and().
		an_error_should_be_logged_with_message("handler for 'cqrs.ConsulEvent' failed: oh no")

}

func Test_ConsulEventDispatcher_dispatch_event_to_single_handler_which_panics_and_second_handler_which_does_not(t *testing.T) {

	given, when, then := ConsulEventDispatcherTest(t)

	given.
		an_event_handler_has_been_registered_which_panics_with("oh no").and().
		a_second_event_handler_has_been_registered()

	when.
		the_event_is_dispatched()

	then.
		the_event_handler_should_be_executed().and().
		the_second_event_handler_should_be_executed().and().
		an_error_should_be_logged_with_message("handler for 'cqrs.ConsulEvent' failed: encountered a panic: oh no")

}

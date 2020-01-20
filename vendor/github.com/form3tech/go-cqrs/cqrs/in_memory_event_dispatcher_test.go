package cqrs

import "testing"

func Test_dispatch_event_to_single_handler(t *testing.T) {

	given, when, then := InMemoryEventDispatcherTest(t)

	given.
		an_event_handler_has_been_registered()

	when.
		the_event_is_dispatched()

	then.
		the_event_handler_should_be_executed().and().
		no_error_should_be_returned()

}

func Test_dispatch_event_to_single_handler_with_context(t *testing.T) {

	given, when, then := InMemoryEventDispatcherTest(t)

	given.
		an_event_handler_with_a_context_has_been_registered()

	when.
		the_event_is_dispatched()

	then.
		the_event_handler_should_be_executed().and().
		no_error_should_be_returned().and().
		the_context_should_be_annotated_with_the_handler()

}

func Test_dispatch_event_multiple_handlers(t *testing.T) {

	given, when, then := InMemoryEventDispatcherTest(t)

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

func Test_register_an_event_handler_that_is_not_a_function(t *testing.T) {

	_, when, then := InMemoryEventDispatcherTest(t)

	when.
		a_handler_that_is_not_a_function_is_registered()

	then.
		a_registration_error_should_be_returned_with_message("handler must be a func with one argument and return error")

}

func Test_register_an_event_handler_with_more_than_one_argument_without_context(t *testing.T) {

	_, when, then := InMemoryEventDispatcherTest(t)

	when.
		a_handler_that_with_more_than_one_argument_is_registered_without_using_context()

	then.
		a_registration_error_should_be_returned_with_message("when using 2 arguments first argument should be of type *context.Context")

}

func Test_register_an_event_handler_with_no_arguments(t *testing.T) {

	_, when, then := InMemoryEventDispatcherTest(t)

	when.
		a_handler_that_with_no_arguments_is_registered()

	then.
		a_registration_error_should_be_returned_with_message("handler must be a func with one argument and return error")

}

func Test_register_an_event_handler_with_more_than_one_return_value(t *testing.T) {

	_, when, then := InMemoryEventDispatcherTest(t)

	when.
		a_handler_that_with_more_than_one_return_value_is_registered()

	then.
		a_registration_error_should_be_returned_with_message("handler must be a func with one argument and return error")

}

func Test_register_an_event_handler_with_no_return_value(t *testing.T) {

	_, when, then := InMemoryEventDispatcherTest(t)

	when.
		a_handler_that_with_no_return_value_is_registered()

	then.
		a_registration_error_should_be_returned_with_message("handler must be a func with one argument and return error")

}

func Test_handler_that_returns_an_error(t *testing.T) {

	given, when, then := InMemoryEventDispatcherTest(t)

	given.
		an_event_handler_has_been_registered_that_will_return_an_error("something went bang")

	when.
		the_event_is_dispatched()

	then.
		an_error_should_be_returned_with_message("error in handler #0 for 'cqrs.testEvent': something went bang")
}

func Test_handler_that_returns_a_panic(t *testing.T) {

	given, when, then := InMemoryEventDispatcherTest(t)

	given.
		an_event_handler_has_been_registered_that_will_panic("oh no: a catastrophic failure occurred")

	when.
		the_event_is_dispatched()

	then.
		an_error_should_be_returned_with_message("error in handler #0 for 'cqrs.testEvent': encountered a panic: oh no: a catastrophic failure occurred")
}

func Test_dispatch_event_multiple_handlers_when_first_handler_errors(t *testing.T) {

	given, when, then := InMemoryEventDispatcherTest(t)

	given.
		an_event_handler_has_been_registered_that_will_return_an_error("something went bang").and().
		a_second_event_handler_has_been_registered()

	when.
		the_event_is_dispatched()

	then.
		the_event_handler_should_be_executed().and().
		the_second_event_handler_should_be_executed().and().
		an_error_should_be_returned_with_message("error in handler #0 for 'cqrs.testEvent': something went bang")
}

func Test_dispatch_event_multiple_handlers_when_first_handler_panics(t *testing.T) {

	given, when, then := InMemoryEventDispatcherTest(t)

	given.
		an_event_handler_has_been_registered_that_will_panic("oh no: a catastrophic failure occurred").and().
		a_second_event_handler_has_been_registered()

	when.
		the_event_is_dispatched()

	then.
		the_event_handler_should_be_executed().and().
		the_second_event_handler_should_be_executed().and().
		an_error_should_be_returned_with_message("error in handler #0 for 'cqrs.testEvent': encountered a panic: oh no: a catastrophic failure occurred")
}

package cqrs

import (
	"testing"

	"github.com/google/uuid"
)

func Test_execute_single_result_query_when_user_has_access(t *testing.T) {

	given, when, then := QueryExecutorTest(t)

	organisationId := uuid.New()

	given.
		a_query_is_registered_with_organisation_filter_that_returns_a_single_result_belonging_to_organisation_id(organisationId)

	when.
		the_query_is_executed_for_a_user_that_has_access_to_organisation_id(organisationId)

	then.
		the_result_is_returned().and().
		no_error_is_returned()

}

func Test_execute_single_result_query_when_user_does_not_have_access(t *testing.T) {

	given, when, then := QueryExecutorTest(t)

	organisationId := uuid.New()
	differentOrganisationId := uuid.New()

	given.
		a_query_is_registered_with_organisation_filter_that_returns_a_single_result_belonging_to_organisation_id(organisationId)

	when.
		the_query_is_executed_for_a_user_that_has_access_to_organisation_id(differentOrganisationId)

	then.
		the_result_is_not_returned().and().
		an_error_is_returned()

}

func Test_execute_multi_result_query_when_user_has_access_to_all_results(t *testing.T) {

	given, when, then := QueryExecutorTest(t)

	organisationId1 := uuid.New()
	organisationId2 := uuid.New()

	given.
		a_query_is_registered_with_organisation_filter_that_returns_returns_for_organisations(organisationId1, organisationId2)

	when.
		the_query_is_executed_by_a_user_that_has_access_to_organisations(organisationId1, organisationId2)

	then.
		all_of_the_results_are_returned().and().
		no_error_is_returned()

}

func Test_execute_multi_result_query_when_user_has_access_to_some_results(t *testing.T) {

	given, when, then := QueryExecutorTest(t)

	organisationId1 := uuid.New()
	organisationId2 := uuid.New()

	given.
		a_query_is_registered_with_organisation_filter_that_returns_returns_for_organisations(organisationId1, organisationId2)

	when.
		the_query_is_executed_by_a_user_that_has_access_to_organisations(organisationId1)

	then.
		only_results_for_the_organisation_the_user_has_access_to_are_returned(organisationId1).and().
		no_error_is_returned()

}

func Test_execute_multi_result_query_when_user_has_access_to_no_results(t *testing.T) {

	given, when, then := QueryExecutorTest(t)

	organisationId1 := uuid.New()
	organisationId2 := uuid.New()
	organisationId3 := uuid.New()

	given.
		a_query_is_registered_with_organisation_filter_that_returns_returns_for_organisations(organisationId1, organisationId2)

	when.
		the_query_is_executed_by_a_user_that_has_access_to_organisations(organisationId3)

	then.
		no_results_are_returned().and().
		no_error_is_returned()
}

func Test_execute_a_db_query_single_result_when_user_has_access(t *testing.T) {

	given, when, then := QueryExecutorTest(t)

	organisationId1 := uuid.New()

	given.
		a_database_query_is_registered_with_organisation_filter_that_returns_a_single_result_belonging_to_organisation_id(organisationId1)

	when.
		the_database_query_is_executed_that_returns_a_result_for_organisation(organisationId1)

	then.
		the_result_is_returned().and().
		no_error_is_returned()
}

func Test_execute_a_single_result_query_with_no_filter(t *testing.T) {

	given, when, then := QueryExecutorTest(t)

	given.
		a_query_is_registered_with_no_result_filter_that_returns_a_single_result()

	when.
		the_query_is_executed_with_no_result_filter_that_returns_a_single_result()

	then.
		the_result_with_no_filter_is_returned().and().
		no_error_is_returned()

}

func Test_execute_a_multi_result_query_with_no_filter(t *testing.T) {

	given, when, then := QueryExecutorTest(t)

	given.
		a_query_is_registered_with_no_result_filter_that_returns_multiple_results()

	when.
		the_query_is_executed_with_no_result_filter_that_returns_multiple_results()

	then.
		all_of_the_no_filter_results_are_returned().and().
		no_error_is_returned()

}

func Test_register_a_query_that_is_not_a_func(t *testing.T) {

	_, when, then := QueryExecutorTest(t)

	when.
		a_non_function_is_registered_with_the_query_executor()

	then.
		an_error_is_returned()
}

func Test_register_a_query_with_more_than_one_argument(t *testing.T) {

	_, when, then := QueryExecutorTest(t)

	when.
		a_function_with_more_than_one_argument_is_registered_with_the_query_executor()

	then.
		an_error_is_returned()
}

func Test_register_a_query_with_only_one_return_value(t *testing.T) {

	_, when, then := QueryExecutorTest(t)

	when.
		a_function_with_only_one_return_value_is_registered_with_the_query_executor()

	then.
		an_error_is_returned()
}

func Test_query_called_without_passing_ptr_to_ptr_for_result(t *testing.T) {

	given, when, then := QueryExecutorTest(t)

	given.
		a_query_is_registered_with_no_result_filter_that_returns_a_single_result()

	when.
		the_query_is_executed_passing_a_ptr_to_result()

	then.
		an_error_is_returned()
}

func Test_register_a_query_with_the_first_return_argument_not_as_a_ptr_to_ptr(t *testing.T) {

	_, when, then := QueryExecutorTest(t)

	when.
		a_function_with_the_first_return_parameter_not_as_a_ptr_is_registered_with_the_query_executor()

	then.
		an_error_is_returned()
}

func Test_register_a_query_with_second_return_argument_is_an_error(t *testing.T) {

	_, when, then := QueryExecutorTest(t)

	when.
		a_function_with_the_second_return_parameter_is_not_an_error_registered_with_the_query_executor()

	then.
		an_error_is_returned()
}

func Test_execute_a_query_that_has_not_been_registered(t *testing.T) {

	_, when, then := QueryExecutorTest(t)

	when.
		a_query_is_executed_that_has_not_been_registered_with_the_query_executor()

	then.
		an_error_is_returned()
}

func Test_result_is_nil_when_query_returns_an_error(t *testing.T) {

	given, when, then := QueryExecutorTest(t)

	given.
		a_query_is_registered_that_will_return_an_error()

	when.
		the_query_that_will_return_an_error_is_executed()

	then.
		the_test_result_is_nil().and().
		an_error_is_returned()
}

func Test_query_returns_nil_single_result(t *testing.T) {

	given, when, then := QueryExecutorTest(t)

	given.
		a_query_is_registered_that_will_return_a_nil_single_result()

	when.
		the_query_is_executed_that_will_return_a_nil_single_result()

	then.
		the_result_is_nil().and().
		no_error_is_returned()

}

func Test_query_returns_nil_mulitple_result(t *testing.T) {

	given, when, then := QueryExecutorTest(t)

	given.
		a_query_is_registered_that_will_return_a_nil_multiple_results()

	when.
		the_query_is_executed_that_will_return_a_nil_multiple_results()

	then.
		the_multiple_result_is_nil().and().
		no_error_is_returned()

}

func Test_query_with_permissions_returns_nil_multiple_result(t *testing.T) {

	given, when, then := QueryExecutorTest(t)

	given.
		a_query_that_uses_context_is_registered_that_will_return_a_nil_multiple_results()

	when.
		the_query_is_executed_that_will_return_a_nil_multiple_results()

	then.
		the_multiple_result_is_nil().and().
		no_error_is_returned()

}

func Test_execute_a_db_query_that_uses_context(t *testing.T) {

	given, when, then := QueryExecutorTest(t)

	organisationId1 := uuid.New()

	given.
		a_database_query_with_context_is_registered_that_returns_a_single_result_belonging_to_organisation_id(organisationId1)

	when.
		the_database_query_is_executed_that_returns_a_result_for_organisation(organisationId1)

	then.
		the_result_is_returned().and().
		no_error_is_returned().and().
		the_context_should_be_annotated_with_the_handler()
}

func Test_execute_a_query_that_uses_context_and_no_db(t *testing.T) {

	given, when, then := QueryExecutorTest(t)

	organisationId1 := uuid.New()

	given.
		a_query_with_context_is_registered_that_returns_a_single_result_belonging_to_organisation_id(organisationId1)

	when.
		the_database_query_is_executed_that_returns_a_result_for_organisation(organisationId1)

	then.
		the_result_is_returned().and().
		no_error_is_returned()
}

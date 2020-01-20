package cqrs

import (
	"context"
	"testing"

	"fmt"

	linq "github.com/ahmetb/go-linq"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type queryExecutorStage struct {
	t             *testing.T
	testResult    *testResult
	testResult2   *testResult2
	multiResults  *[]testResult
	multiResults2 *[]testResult2
	err           error
	queryExecutor QueryExecutor
	annotation    string
}

type testCriteria struct {
	a int
	b int
}

type testResult struct {
	answer         int
	OrganisationId uuid.UUID
}

type testCriteria2 struct {
	foo int
}

type testResult2 struct {
	bar int
}

func authErrorFunc(s string) error {
	return fmt.Errorf(s)
}

func QueryExecutorTest(t *testing.T) (*queryExecutorStage, *queryExecutorStage, *queryExecutorStage) {
	stage := &queryExecutorStage{
		t: t,
		queryExecutor: &queryExecutor{
			queries: make(map[string]map[string]queryHolder),
			db:      &sqlx.DB{},
		},
	}

	return stage, stage, stage
}

func (s *queryExecutorStage) a_query_is_registered_with_organisation_filter_that_returns_a_single_result_belonging_to_organisation_id(organisationId uuid.UUID) *queryExecutorStage {

	err := s.queryExecutor.RegisterQuery(func(criteria testCriteria) (*testResult, error) {

		return &testResult{answer: criteria.a + criteria.b, OrganisationId: organisationId}, nil
	}, WithOrganisationFilter(authErrorFunc, func(ctx *context.Context) ([]uuid.UUID, error) {

		c := *ctx

		return c.Value("organisationIds").([]uuid.UUID), nil
	}))
	if err != nil {
		s.t.Fatalf("error registering query, error: %v", err)
	}

	return s
}

func (s *queryExecutorStage) the_query_is_executed_for_a_user_that_has_access_to_organisation_id(organisationId uuid.UUID) *queryExecutorStage {

	s.testResult = &testResult{}

	ctx := context.WithValue(context.Background(), "organisationIds", []uuid.UUID{organisationId})

	s.err = s.queryExecutor.Execute(&ctx, testCriteria{a: 1, b: 2}, &s.testResult)

	return s
}

func (s *queryExecutorStage) the_query_is_executed_by_a_user_that_has_access_to_organisations(organisationIds ...uuid.UUID) *queryExecutorStage {

	s.multiResults = &[]testResult{}

	ctx := context.WithValue(context.Background(), "organisationIds", organisationIds)

	s.err = s.queryExecutor.Execute(&ctx, testCriteria{}, &s.multiResults)

	return s
}

func (s *queryExecutorStage) a_function_with_more_than_one_argument_is_registered_with_the_query_executor() *queryExecutorStage {

	s.err = s.queryExecutor.RegisterQuery(func(a, b int) (*testResult, error) { return &testResult{answer: a + b}, nil }, nil)

	return s
}

func (s *queryExecutorStage) a_query_is_executed_that_has_not_been_registered_with_the_query_executor() *queryExecutorStage {

	s.testResult = &testResult{}
	s.err = s.queryExecutor.Execute(nil, testCriteria{a: 1, b: 2}, &s.testResult)

	return s
}

func (s *queryExecutorStage) no_results_are_returned() *queryExecutorStage {

	assert.Equal(s.t, 0, len(*s.multiResults))

	return s
}

func (s *queryExecutorStage) a_query_is_registered_with_no_result_filter_that_returns_a_single_result() *queryExecutorStage {

	err := s.queryExecutor.RegisterQuery(func(criteria testCriteria2) (*testResult2, error) {

		return &testResult2{bar: criteria.foo}, nil
	}, WithNoFilter())
	if err != nil {
		s.t.Fatalf("error registering query, error: %v", err)
	}

	return s
}

func (s *queryExecutorStage) a_query_is_registered_that_will_return_a_nil_single_result() *queryExecutorStage {

	err := s.queryExecutor.RegisterQuery(func(criteria testCriteria2) (*testResult2, error) {

		return nil, nil
	}, WithNoFilter())
	if err != nil {
		s.t.Fatalf("error registering query, error: %v", err)
	}

	return s
}

func (s *queryExecutorStage) a_query_is_registered_that_will_return_a_nil_multiple_results() *queryExecutorStage {

	err := s.queryExecutor.RegisterQuery(func(criteria testCriteria2) (*[]testResult2, error) {

		return nil, nil
	}, WithNoFilter())
	if err != nil {
		s.t.Fatalf("error registering query, error: %v", err)
	}

	return s
}

func (s *queryExecutorStage) a_query_that_uses_context_is_registered_that_will_return_a_nil_multiple_results() *queryExecutorStage {

	err := s.queryExecutor.RegisterQuery(func(ctx *context.Context, criteria testCriteria2) (*[]testResult2, error) {

		return nil, nil
	}, WithNoFilter())
	if err != nil {
		s.t.Fatalf("error registering query, error: %v", err)
	}

	return s
}

func (s *queryExecutorStage) the_query_is_executed_with_no_result_filter_that_returns_a_single_result() *queryExecutorStage {

	s.testResult2 = &testResult2{}

	s.err = s.queryExecutor.Execute(nil, testCriteria2{foo: 123}, &s.testResult2)

	return s
}

func (s *queryExecutorStage) the_query_is_executed_passing_a_ptr_to_result() *queryExecutorStage {

	s.testResult2 = &testResult2{}

	s.err = s.queryExecutor.Execute(nil, testCriteria2{foo: 123}, s.testResult2)

	return s
}

func (s *queryExecutorStage) the_query_is_executed_that_will_return_a_nil_single_result() *queryExecutorStage {

	s.testResult2 = &testResult2{}

	s.err = s.queryExecutor.Execute(nil, testCriteria2{foo: 123}, &s.testResult2)

	return s
}

func (s *queryExecutorStage) the_query_is_executed_that_will_return_a_nil_multiple_results() *queryExecutorStage {

	s.multiResults2 = &[]testResult2{}

	s.err = s.queryExecutor.Execute(nil, testCriteria2{foo: 123}, &s.multiResults2)

	return s
}

func (s *queryExecutorStage) a_query_is_registered_that_will_return_an_error() *queryExecutorStage {

	err := s.queryExecutor.RegisterQuery(func(criteria testCriteria) (*testResult, error) {

		return nil, fmt.Errorf("error")
	}, WithOrganisationFilter(authErrorFunc, func(ctx *context.Context) ([]uuid.UUID, error) {

		c := *ctx

		return c.Value("organisationIds").([]uuid.UUID), nil
	}))
	if err != nil {
		s.t.Fatalf("error registering query, error: %v", err)
	}

	return s
}

func (s *queryExecutorStage) the_query_that_will_return_an_error_is_executed() *queryExecutorStage {

	s.testResult = &testResult{}
	s.err = s.queryExecutor.Execute(nil, testCriteria{a: 1, b: 2}, &s.testResult)

	return s
}

func (s *queryExecutorStage) a_query_is_registered_with_organisation_filter_that_returns_returns_for_organisations(organisationIds ...uuid.UUID) *queryExecutorStage {

	err := s.queryExecutor.RegisterQuery(func(criteria testCriteria) (*[]testResult, error) {

		var results []testResult

		for i, organisationId := range organisationIds {
			results = append(results, testResult{answer: i, OrganisationId: organisationId})
		}

		return &results, nil
	}, WithOrganisationFilter(authErrorFunc, func(ctx *context.Context) ([]uuid.UUID, error) {

		c := *ctx

		return c.Value("organisationIds").([]uuid.UUID), nil
	}))
	if err != nil {
		s.t.Fatalf("error registering query, error: %v", err)
	}

	return s
}

func (s *queryExecutorStage) a_non_function_is_registered_with_the_query_executor() *queryExecutorStage {
	s.err = s.queryExecutor.RegisterQuery(1, WithOrganisationFilter(authErrorFunc, func(ctx *context.Context) ([]uuid.UUID, error) {
		return nil, nil
	}))

	return s
}

func (s *queryExecutorStage) a_function_with_only_one_return_value_is_registered_with_the_query_executor() *queryExecutorStage {

	s.err = s.queryExecutor.RegisterQuery(func(criteria testCriteria) error { return nil }, nil)

	return s
}

func (s *queryExecutorStage) a_function_with_the_first_return_parameter_not_as_a_ptr_is_registered_with_the_query_executor() *queryExecutorStage {

	s.err = s.queryExecutor.RegisterQuery(func(criteria testCriteria) (testResult, error) { return testResult{}, nil }, nil)

	return s
}

func (s *queryExecutorStage) a_function_with_the_second_return_parameter_is_not_an_error_registered_with_the_query_executor() *queryExecutorStage {

	s.err = s.queryExecutor.RegisterQuery(func(criteria testCriteria) (testResult, testResult) { return testResult{}, testResult{} }, nil)

	return s
}

func (s *queryExecutorStage) a_database_query_is_registered_with_organisation_filter_that_returns_a_single_result_belonging_to_organisation_id(organisationId uuid.UUID) *queryExecutorStage {

	err := s.queryExecutor.RegisterQuery(func(db *sqlx.DB, criteria testCriteria) (*testResult, error) {

		return &testResult{answer: criteria.a + criteria.b, OrganisationId: organisationId}, nil
	}, WithOrganisationFilter(authErrorFunc, func(ctx *context.Context) ([]uuid.UUID, error) {

		c := *ctx

		return c.Value("organisationIds").([]uuid.UUID), nil
	}))
	if err != nil {
		s.t.Fatalf("error registering query, error: %v", err)
	}

	return s
}

func (s *queryExecutorStage) the_database_query_is_executed_that_returns_a_result_for_organisation(organisationId uuid.UUID) *queryExecutorStage {

	s.testResult = &testResult{}

	ctx := context.WithValue(context.Background(), "organisationIds", []uuid.UUID{organisationId})

	s.err = s.queryExecutor.Execute(&ctx, testCriteria{a: 1, b: 2}, &s.testResult)

	return s
}

func (s *queryExecutorStage) the_result_is_returned() *queryExecutorStage {

	assert.Equal(s.t, 3, s.testResult.answer)

	return s
}

func (s *queryExecutorStage) a_query_is_registered_with_no_result_filter_that_returns_multiple_results() *queryExecutorStage {

	err := s.queryExecutor.RegisterQuery(func(criteria testCriteria2) (*[]testResult2, error) {

		var results []testResult2

		for i := 0; i < 5; i++ {
			results = append(results, testResult2{bar: i})
		}

		return &results, nil
	}, WithNoFilter())
	if err != nil {
		s.t.Fatalf("error registering query, error: %v", err)
	}
	return s
}

func (s *queryExecutorStage) the_query_is_executed_with_no_result_filter_that_returns_multiple_results() *queryExecutorStage {

	s.multiResults2 = &[]testResult2{}
	s.err = s.queryExecutor.Execute(nil, testCriteria2{foo: 1}, &s.multiResults2)

	return s
}

func (s *queryExecutorStage) and() *queryExecutorStage {
	return s
}

func (s *queryExecutorStage) no_error_is_returned() *queryExecutorStage {

	assert.Nil(s.t, s.err)

	return s
}

func (s *queryExecutorStage) the_result_is_not_returned() *queryExecutorStage {

	assert.Equal(s.t, 0, s.testResult.answer)

	return s
}

func (s *queryExecutorStage) an_error_is_returned() *queryExecutorStage {

	assert.NotNil(s.t, s.err)

	return s
}

func (s *queryExecutorStage) all_of_the_results_are_returned() *queryExecutorStage {

	assert.Equal(s.t, 2, len(*s.multiResults))

	return s
}

func (s *queryExecutorStage) only_results_for_the_organisation_the_user_has_access_to_are_returned(organisationId uuid.UUID) *queryExecutorStage {

	assert.Equal(s.t, 1, len(*s.multiResults))
	assert.Equal(s.t, 0, linq.From(*s.multiResults).WhereT(func(testResult testResult) bool { return testResult.OrganisationId != organisationId }).Count())

	return s
}
func (s *queryExecutorStage) the_result_with_no_filter_is_returned() *queryExecutorStage {

	assert.Equal(s.t, 123, s.testResult2.bar)

	return s
}

func (s *queryExecutorStage) all_of_the_no_filter_results_are_returned() *queryExecutorStage {

	assert.Equal(s.t, 5, len(*s.multiResults2))

	return s
}

func (s *queryExecutorStage) the_test_result_is_nil() *queryExecutorStage {

	assert.Nil(s.t, s.testResult)

	return s
}

func (s *queryExecutorStage) the_result_is_nil() *queryExecutorStage {

	assert.Nil(s.t, s.testResult2)

	return s
}

func (s *queryExecutorStage) the_multiple_result_is_nil() *queryExecutorStage {

	assert.Nil(s.t, s.multiResults2)

	return s
}

func (s *queryExecutorStage) a_database_query_with_context_is_registered_that_returns_a_single_result_belonging_to_organisation_id(organisationId uuid.UUID) *queryExecutorStage {

	err := s.queryExecutor.RegisterQuery(func(ctx *context.Context, db *sqlx.DB, criteria testCriteria) (*testResult, error) {
		if ctx != nil {
			s.annotation = (*ctx).Value("handler").(string)
		}
		return &testResult{answer: criteria.a + criteria.b, OrganisationId: organisationId}, nil
	}, WithNoFilter())
	if err != nil {
		s.t.Fatalf("error registering query, error: %v", err)
	}

	return s
}

func (s *queryExecutorStage) a_query_with_context_is_registered_that_returns_a_single_result_belonging_to_organisation_id(organisationId uuid.UUID) *queryExecutorStage {

	err := s.queryExecutor.RegisterQuery(func(ctx *context.Context, criteria testCriteria) (*testResult, error) {

		return &testResult{answer: criteria.a + criteria.b, OrganisationId: organisationId}, nil
	}, WithNoFilter())
	if err != nil {
		s.t.Fatalf("error registering query, error: %v", err)
	}

	return s
}

func (s *queryExecutorStage) the_context_should_be_annotated_with_the_handler() *queryExecutorStage {
	assert.Equal(s.t, "github.com/form3tech/go-cqrs/cqrs.(*queryExecutorStage).a_database_query_with_context_is_registered_that_returns_a_single_result_belonging_to_organisation_id.func1", s.annotation)
	return s
}

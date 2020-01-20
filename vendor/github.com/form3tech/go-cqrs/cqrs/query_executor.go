package cqrs

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"sync"

	linq "github.com/ahmetb/go-linq"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type QueryExecutor interface {
	Execute(ctx *context.Context, criteria interface{}, result interface{}) error
	RegisterQuery(query interface{}, filter queryFilter) error
}

var oneQueryExecutor sync.Once
var queryExecutorInstance *queryExecutor

const (
	errorString = " func(c criteria) (result, error) or func(db *sqlx.db. c criteria) (*result, error)"
)

type queryHolder struct {
	query       interface{}
	usesDb      bool
	usesContext bool
	queryFilter queryFilter
}

type queryExecutor struct {
	db      *sqlx.DB
	queries map[string]map[string]queryHolder
}

type queryFilter func(ctx *context.Context, queryResult []reflect.Value, result interface{}) error

func GetQueryExecutor(db *sqlx.DB) QueryExecutor {

	oneQueryExecutor.Do(func() {
		queryExecutorInstance = &queryExecutor{
			queries: make(map[string]map[string]queryHolder),
			db:      db,
		}
	})

	return queryExecutorInstance
}

func WithNoFilter() queryFilter {
	return func(ctx *context.Context, queryResult []reflect.Value, result interface{}) error {

		v := reflect.ValueOf(result)
		p := v.Elem()
		p.Set(queryResult[0])

		return nil
	}
}

func WithOrganisationFilter(authError func(string) error, getAllowedOrganisationIds func(c *context.Context) ([]uuid.UUID, error)) queryFilter {
	return func(ctx *context.Context, queryResult []reflect.Value, result interface{}) error {
		allowedOrganisationIds, err := getAllowedOrganisationIds(ctx)

		if err != nil {
			return err
		}

		resultType := reflect.TypeOf(result)
		resultKind := resultType.Elem().Elem().Kind()
		if resultKind == reflect.Array || resultKind == reflect.Slice {

			multiResults := reflect.New(resultType.Elem().Elem()).Elem()
			resultValue := reflect.ValueOf(queryResult[0].Elem().Interface())

			for i := 0; i < resultValue.Len(); i++ {

				organisationIdValue := reflect.Indirect(resultValue.Index(i)).FieldByName("OrganisationId")
				organisationId := organisationIdValue.Interface().(uuid.UUID)

				count := linq.From(allowedOrganisationIds).
					WhereT(func(o uuid.UUID) bool { return o == organisationId }).
					Count()

				if count > 0 {
					multiResults.Set(reflect.Append(multiResults, resultValue.Index(i)))
				}
			}

			v := reflect.ValueOf(result).Elem().Elem()
			v.Set(multiResults)

		} else {

			organisationIdValue := reflect.Indirect(queryResult[0]).FieldByName("OrganisationId")
			organisationId := organisationIdValue.Interface().(uuid.UUID)

			count := linq.From(allowedOrganisationIds).
				WhereT(func(o uuid.UUID) bool { return o == organisationId }).
				Count()

			if count == 0 {
				return authError(fmt.Sprintf("user does not have access to result for organisation id: %v", organisationId))
			}

			if err != nil {
				return err
			}

			v := reflect.ValueOf(result)
			p := v.Elem()
			p.Set(queryResult[0])

		}

		return nil
	}
}

/*

When calling execute you must pass **result for your result even though your query returns *result.

This is to handle the fact that your query may return nil and you cannot set the top level ptr if you just pass *result

*/

func (q *queryExecutor) Execute(ctx *context.Context, criteria interface{}, result interface{}) error {

	criteriaType := reflect.TypeOf(criteria)
	queries, ok := q.queries[criteriaType.String()]

	if !ok {
		return fmt.Errorf("no query found taking parameter of type %v, did you forget to register the query?  Call RegisterQuery to register", criteriaType)
	}

	resultType := reflect.TypeOf(result)

	if resultType.Kind() != reflect.Ptr {
		return fmt.Errorf("you must pass pointer to pointer your varable for result parameter (this is where the query result will be stored)")
	}

	if resultType.Elem().Kind() != reflect.Ptr {
		return fmt.Errorf("you must pass pointer to pointer your varable for result parameter (this is where the query result will be stored)")
	}

	underlyingResultType := resultType.Elem().String()
	query, ok := queries[underlyingResultType]

	if !ok {
		return fmt.Errorf("found queries accepting criteria %s but none returning result %s", criteriaType.String(), resultType.String())
	}

	method := reflect.ValueOf(query.query)

	if ctx != nil {
		annotatedContext := context.WithValue(*ctx, "handler", runtime.FuncForPC(method.Pointer()).Name())
		ctx = &annotatedContext
	}

	var queryResult []reflect.Value

	if query.usesDb && !query.usesContext {
		queryResult = method.Call([]reflect.Value{reflect.ValueOf(q.db), reflect.ValueOf(criteria)})
	} else if !query.usesDb && query.usesContext {
		queryResult = method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(criteria)})
	} else if query.usesDb && query.usesContext {
		queryResult = method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(q.db), reflect.ValueOf(criteria)})
	} else {
		queryResult = method.Call([]reflect.Value{reflect.ValueOf(criteria)})
	}

	if queryResult[1].Interface() != nil {

		v := reflect.ValueOf(result)
		p := v.Elem()
		p.Set(reflect.Zero(p.Type()))

		return queryResult[1].Interface().(error)
	}

	if queryResult[0].IsNil() {

		v := reflect.ValueOf(result)
		p := v.Elem()
		p.Set(reflect.Zero(p.Type()))

		return nil
	}

	return query.queryFilter(ctx, queryResult, result)

}

func (q *queryExecutor) RegisterQuery(query interface{}, filter queryFilter) error {
	queryType := reflect.TypeOf(query)
	usesDb, usesContext := false, false
	var dbType *sqlx.DB
	var ctxType *context.Context
	var criteriaType string

	if !(queryType.Kind() == reflect.Func) {
		return fmt.Errorf("%s is not a func must be a func", queryType.Kind())
	}

	if queryType.NumIn() < 1 || queryType.NumIn() > 3 {
		return fmt.Errorf("%s must only have one argument to three arguments signature should accept (criteria), (ctx, criteria) or (ctx, *context.Context, *sqlx.Db, criteria)", queryType.Kind())
	}

	if queryType.NumIn() == 1 {
		criteriaType = queryType.In(0).String()
	}

	if queryType.NumIn() == 2 {
		if queryType.In(0) == reflect.TypeOf(dbType) {
			usesDb = true
		} else if queryType.In(0) == reflect.TypeOf(ctxType) {
			usesContext = true
		} else {
			return fmt.Errorf("when using 2 arguments first argument should be of type *sqlx.db or *context.Context")
		}

		criteriaType = queryType.In(1).String()
	}

	if queryType.NumIn() == 3 {
		if queryType.In(0) == reflect.TypeOf(ctxType) {
			usesContext = true
		} else {
			return fmt.Errorf("when using 3 arguments first argument should be of type *context.Context")
		}

		if queryType.In(1) == reflect.TypeOf(dbType) {
			usesDb = true
		} else {
			return fmt.Errorf("when using 3 arguments second argument should be of type *sqlx.Db")
		}

		criteriaType = queryType.In(2).String()
	}

	if queryType.NumOut() != 2 {
		return fmt.Errorf("%s must return 2 parameters (result, error)", queryType.Kind())
	}

	if queryType.Out(0).Kind() != reflect.Ptr {
		return fmt.Errorf("first return parameter of query func must be a pointer to your query result signature should be %s", errorString)
	}

	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if !queryType.Out(1).Implements(errorInterface) {
		return fmt.Errorf("first return parameter of query func must be error signature should be %s", errorString)
	}

	resultType := queryType.Out(0).String()

	m, ok := q.queries[criteriaType]

	if !ok {
		m = make(map[string]queryHolder)
		q.queries[criteriaType] = m
	}

	if _, ok := m[resultType]; ok {
		return fmt.Errorf("query already registered with criteria type: %v and result type: %v", queryType, resultType)
	}

	m[resultType] = queryHolder{
		query:       query,
		queryFilter: filter,
		usesDb:      usesDb,
		usesContext: usesContext,
	}

	return nil
}

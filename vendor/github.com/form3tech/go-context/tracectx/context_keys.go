package tracectx

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
)

// kept private to force usage of NewContextKey function
type contextKey struct {
	// unique key name, such as correlation_id or transaction_time
	Key string
	// an Http Header representation for the this context Key, such as XForm3-Correlation-Id
	HTTPHeader string
	// Converters, for transforming the value stored in a context.Context to/from a string from/to typed value
	Conversions BidirectionalConversion
}

type ContextKey interface {
	HttpHeader() string
	String() string
	ParseValue(val string) (interface{}, error)
	StringFromContext(ctx context.Context) (string, error)
}

var (
	CorrelationId, _ = NewStringContextKey("correlation_id", "X-Form3-correlation-id")
	// start time of the transaction
	TransactionTime, _ = NewTimeContextKey("transaction_time", "X-Form3-Transaction-Time")
	// duration of a transaction, calculated in place based on the TransactionTime and current time
	TransactionDuration, _ = NewDurationContextKey("transaction_duration", "X-Form3-Transaction-Duration")
	// Time that should be excluded from calculations of any sort (for instance, performance metrics)
	TransactionDurationExternal, _ = NewDurationContextKey("transaction_duration_external", "X-Form3-Transaction-Duration-External")
	Handler, _                     = NewTimeContextKey("handler", "")

	availableContextKeys = make(map[string]ContextKey)
)

func NewStringContextKey(key, httpHeader string) (ContextKey, error) {
	return mustNewContextKey(contextKey{
		Key:         key,
		HTTPHeader:  httpHeader,
		Conversions: StringConverter,
	})
}

func NewTimeContextKey(key, httpHeader string) (ContextKey, error) {
	return mustNewContextKey(contextKey{
		Key:         key,
		HTTPHeader:  httpHeader,
		Conversions: TimeConverter,
	})
}

func NewDurationContextKey(key, httpHeader string) (ContextKey, error) {
	return mustNewContextKey(contextKey{
		Key:         key,
		HTTPHeader:  httpHeader,
		Conversions: DurationConverter,
	})
}

func NewIntContextKey(key, httpHeader string) (ContextKey, error) {
	return mustNewContextKey(contextKey{
		Key:         key,
		HTTPHeader:  httpHeader,
		Conversions: IntConverter,
	})
}

func NewInt32ContextKey(key, httpHeader string) (ContextKey, error) {
	return mustNewContextKey(contextKey{
		Key:         key,
		HTTPHeader:  httpHeader,
		Conversions: Int32Converter,
	})
}

func NewInt64ContextKey(key, httpHeader string) (ContextKey, error) {
	return mustNewContextKey(contextKey{
		Key:         key,
		HTTPHeader:  httpHeader,
		Conversions: Int64Converter,
	})
}

// GetContextKey returns a context key, if available, associated to a string key name, such as correlation_id or transaction_time.
// If this returns null, functions such as NewStringContextKey, NewDurationContextKey or NewTimeContextKey can be used to register
// a new ContextKey instance
func GetContextKey(key string) ContextKey {
	if ctxKey, exists := availableContextKeys[key]; exists {
		return ctxKey
	}
	return nil
}

func (k *contextKey) HttpHeader() string {
	return k.HTTPHeader
}

func (k *contextKey) String() string {
	return k.Key
}

func (k *contextKey) ParseValue(val string) (interface{}, error) {
	return k.Conversions.Parser(val)
}

func (k *contextKey) StringFromContext(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("context is nil")
	}
	ctxVal := ctx.Value(k)
	if ctxVal != nil {
		return k.Conversions.Stringer(ctxVal)
	}
	return "", fmt.Errorf("no value in context for %s", k.String())
}

// mustNewContextKey ensures a new context key is created with all required fields, as well as adding it to the
// map of existing context keys
func mustNewContextKey(key contextKey) (ContextKey, error) {
	err := validation.ValidateStruct(&key,
		validation.Field(&key.Conversions, validation.Required),
		validation.Field(&key.Key, validation.Required),
	)
	if err != nil {
		return nil, err
	}

	availableContextKeys[key.Key] = &key
	return &key, nil
}

package tracectx

import (
	"context"
	"fmt"
)

// NewContextWithValue stores a value into context.Context, first parsing value to a typed instance, depending on which
// ContextKey was passed in as parameter. ctxKey can be determined by using GetContextKey method. Even if NewContextWithValue
// errors, the original context.Context is returned.
func NewContextWithValue(ctx context.Context, ctxKey ContextKey, value string) (context.Context, error) {
	if ctxKey == nil {
		return ctx, fmt.Errorf("ctxKey passed is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	parsedValue, err := ctxKey.ParseValue(value)
	if err != nil {
		return ctx, err
	}
	if parsedValue != nil {
		return context.WithValue(ctx, ctxKey, parsedValue), nil
	}
	return ctx, nil
}

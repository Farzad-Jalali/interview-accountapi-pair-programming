package internal

import (
	"fmt"
	"strconv"
	"time"
)

var errNil = fmt.Errorf("argument is nil")

// Stringer tries casting a val to string or calling String() method if available
func Stringer(val interface{}) (string, error) {
	if val == nil {
		return "", errNil
	}
	if s, ok := val.(string); ok {
		return s, nil
	}
	if s, ok := val.(fmt.Stringer); ok {
		return s.String(), nil
	}
	return "", fmt.Errorf("%v is neither a string nor has a method String()", val)
}

func TimeStringer(val interface{}) (string, error) {
	if val == nil {
		return "", errNil
	}
	if castTime, ok := val.(time.Time); ok {
		formatted := castTime.Format(time.RFC3339Nano)
		return formatted, nil
	}
	return "", fmt.Errorf("%v not of type time.Time", val)
}

func DurationStringer(val interface{}) (string, error) {
	if val == nil {
		return "", errNil
	}
	if dur, ok := val.(time.Duration); ok {
		return dur.String(), nil
	}
	return "", fmt.Errorf("%v not of type duration.Duration", val)
}

func IntStringer(val interface{}) (string, error) {
	if val == nil {
		return "", errNil
	}
	if castInt, ok := val.(int); ok {
		return strconv.Itoa(castInt), nil
	}
	return "", fmt.Errorf("%v not of type int", val)
}

func Int32Stringer(val interface{}) (string, error) {
	if val == nil {
		return "", errNil
	}
	if castInt, ok := val.(int32); ok {
		return strconv.Itoa(int(castInt)), nil
	}
	return "", fmt.Errorf("%v not of type int32", val)
}

func Int64Stringer(val interface{}) (string, error) {
	if val == nil {
		return "", errNil
	}
	if castInt, ok := val.(int64); ok {
		return strconv.FormatInt(castInt, 64), nil
	}
	return "", fmt.Errorf("%v not of type int64", val)
}

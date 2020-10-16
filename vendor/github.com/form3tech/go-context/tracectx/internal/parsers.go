package internal

import (
	"strconv"
	"time"
)

func TimeParser(value string) (interface{}, error) {
	return time.Parse(time.RFC3339Nano, value)
}

// StringParser simply returns input string value
func StringParser(value string) (interface{}, error) {
	return value, nil
}

func DurationParser(value string) (interface{}, error) {
	return time.ParseDuration(value)
}

func IntParser(value string) (interface{}, error) {
	return strconv.Atoi(value)
}

func Int32Parser(value string) (interface{}, error) {
	parsed, err := strconv.ParseInt(value, 10, 32)
	return int32(parsed), err
}

func Int64Parser(value string) (interface{}, error) {
	return strconv.ParseInt(value, 10, 64)
}

package log

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithFields(fields logrus.Fields) Logger
	WithContext(c *context.Context) Logger
	Data() logrus.Fields
	Time() time.Time
	Level() logrus.Level
	Logger() *logrus.Logger
	Message() string
	Info(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type entryWrapper struct {
	logEntry *logrus.Entry
}

func newEntry(logEntry *logrus.Entry) Logger {
	return &entryWrapper{
		logEntry: logEntry,
	}
}

func (e *entryWrapper) WithContext(c *context.Context) Logger {
	if c == nil {
		return e
	}
	corrId, ok := (*c).Value("correlation-id").(string)
	if ok && corrId != "" {
		return e.WithField("correlation-id", corrId)
	}
	return e
}

func (e *entryWrapper) WithField(key string, value interface{}) Logger {
	return e.WithFields(logrus.Fields{key: value})
}
func (e *entryWrapper) WithFields(fields logrus.Fields) Logger {
	return newEntry(e.logEntry.WithFields(fields))
}
func (e *entryWrapper) Data() logrus.Fields {
	return e.logEntry.Data
}
func (e *entryWrapper) Time() time.Time {
	return e.logEntry.Time
}
func (e *entryWrapper) Level() logrus.Level {
	return e.logEntry.Level
}
func (e *entryWrapper) Logger() *logrus.Logger {
	return e.logEntry.Logger
}
func (e *entryWrapper) Message() string {
	return e.logEntry.Message
}

func (e *entryWrapper) Info(args ...interface{})  { e.logEntry.Info(args...) }
func (e *entryWrapper) Debug(args ...interface{}) { e.logEntry.Debug(args...) }
func (e *entryWrapper) Warn(args ...interface{})  { e.logEntry.Warn(args...) }
func (e *entryWrapper) Error(args ...interface{}) { e.logEntry.Error(args...) }
func (e *entryWrapper) Panic(args ...interface{}) { e.logEntry.Panic(args...) }
func (e *entryWrapper) Fatal(args ...interface{}) { e.logEntry.Fatal(args...) }
func (e *entryWrapper) Debugf(format string, args ...interface{}) {
	e.logEntry.Debugf(format, args...)
}
func (e *entryWrapper) Infof(format string, args ...interface{}) {
	e.logEntry.Infof(format, args...)
}
func (e *entryWrapper) Warnf(format string, args ...interface{}) {
	e.logEntry.Warnf(format, args...)
}
func (e *entryWrapper) Errorf(format string, args ...interface{}) {
	e.logEntry.Errorf(format, args...)
}
func (e *entryWrapper) Panicf(format string, args ...interface{}) {
	e.logEntry.Panicf(format, args...)
}
func (e *entryWrapper) Fatalf(format string, args ...interface{}) {
	e.logEntry.Fatalf(format, args...)
}

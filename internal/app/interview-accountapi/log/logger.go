package log

import (
	"context"
	"strings"

	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/settings"

	"github.com/sirupsen/logrus"
)

const (
	TimeFormat = "2006-01-02T15:04:05.000Z07:00"
)

var log Logger

func init() {
	log = CreateLogger()
}

func CreateLogger() Logger {
	if strings.EqualFold(settings.LogFormat, "json") {
		logrus.SetFormatter(&logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyMsg:  "message",
				logrus.FieldKeyTime: "@timestamp",
			},
			TimestampFormat: TimeFormat,
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyMsg:  "message",
				logrus.FieldKeyTime: "@timestamp",
			},
			TimestampFormat: TimeFormat,
		})
	}

	if level, err := logrus.ParseLevel(settings.LogLevel); err == nil {
		logrus.SetLevel(level)
	}
	log := logrus.
		WithFields(logrus.Fields{
			"stack":        settings.StackName,
			"service_name": settings.ServiceName,
		})
	return newEntry(log)
}

func AddHook(hook logrus.Hook) {
	log.Logger().AddHook(hook)
}

func RemoveHook(hook logrus.Hook) {
	for _, level := range hook.Levels() {
		var newHookAr []logrus.Hook

		for _, h := range log.Logger().Hooks[level] {
			if h != hook {
				newHookAr = append(newHookAr, hook)
			}
		}
		log.Logger().Hooks[level] = newHookAr
	}
}

func StandardLogger() *logrus.Logger {
	return log.Logger()
}

func WithFields(fields map[string]interface{}) Logger {
	return log.WithFields(logrus.Fields(fields))
}
func WithField(key string, value interface{}) Logger {
	return log.WithField(key, value)
}
func WithContext(c *context.Context) Logger {
	return log.WithContext(c)
}

func Debug(args ...interface{})                 { log.Debug(args...) }
func Info(args ...interface{})                  { log.Info(args...) }
func Warn(args ...interface{})                  { log.Warn(args...) }
func Error(args ...interface{})                 { log.Error(args...) }
func Panic(args ...interface{})                 { log.Panic(args...) }
func Fatal(args ...interface{})                 { log.Fatal(args...) }
func Debugf(format string, args ...interface{}) { log.Debugf(format, args...) }
func Infof(format string, args ...interface{})  { log.Infof(format, args...) }
func Warnf(format string, args ...interface{})  { log.Warnf(format, args...) }
func Errorf(format string, args ...interface{}) { log.Errorf(format, args...) }
func Panicf(format string, args ...interface{}) { log.Panicf(format, args...) }
func Fatalf(format string, args ...interface{}) { log.Fatalf(format, args...) }

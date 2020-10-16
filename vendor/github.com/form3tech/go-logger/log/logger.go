package log

import (
	"context"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	TimeFormat         = "2006-01-02T15:04:05.000Z07:00"
	AlertAudienceLabel = "alert_audience"
	envServiceTags     = "SERVICE_TAGS"
)

var versionRegex = regexp.MustCompile(`{version:"(?P<version>[^"]+)"}`)

var log = CreateLogger()

func CreateLogger() Logger {
	logFormat := os.Getenv("LOG_FORMAT")
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	stackName := os.Getenv("STACK_NAME")
	serviceName := os.Getenv("SERVICE_NAME")

	if strings.EqualFold(logFormat, "json") {
		logrus.SetFormatter(NewTransactionTimeFormatter(
			&logrus.JSONFormatter{
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyMsg:  "message",
					logrus.FieldKeyTime: "@timestamp",
				},
				TimestampFormat: TimeFormat,
			}))
	} else {
		logrus.SetFormatter(NewTransactionTimeFormatter(
			&logrus.TextFormatter{
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyMsg:  "message",
					logrus.FieldKeyTime: "@timestamp",
				},
				TimestampFormat: TimeFormat,
			}))
	}

	if level, err := logrus.ParseLevel(logLevel); err == nil {
		logrus.SetLevel(level)
	}

	logFields := logrus.Fields{
		"stack":        stackName,
		"service_name": serviceName,
	}

	if version := withVersionFromServiceTags(); version != "" {
		logFields["service_version"] = version
	}

	return newEntry(logrus.WithFields(logFields))
}

func AddHook(hook logrus.Hook) {
	log.Logger().AddHook(hook)
}

func RemoveHook(hook logrus.Hook) {
	newHooks := make(logrus.LevelHooks)
	for _, level := range logrus.AllLevels {
		var newHookAr []logrus.Hook

		for _, h := range log.Logger().Hooks[level] {
			if h != hook {
				newHookAr = append(newHookAr, h)
			}
		}
		newHooks[level] = newHookAr
	}

	log.Logger().ReplaceHooks(newHooks)
}

func withVersionFromServiceTags() string {
	version := versionRegex.FindStringSubmatch(os.Getenv(envServiceTags))
	if len(version) != 2 {
		return ""
	}

	return version[1]
}

func StandardLogger() *logrus.Logger {
	return log.Logger()
}

func WithFields(fields map[string]interface{}) Logger {
	return log.WithFields(fields)
}
func WithField(key string, value interface{}) Logger {
	return log.WithField(key, value)
}
func WithError(value error) Logger {
	return log.WithError(value)
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
func Alert(audience string, args ...interface{}) {
	log.WithField(AlertAudienceLabel, audience).Error(args...)
}
func Alertf(audience string, format string, args ...interface{}) {
	log.WithField(AlertAudienceLabel, audience).Errorf(format, args...)
}

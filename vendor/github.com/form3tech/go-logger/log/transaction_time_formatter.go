package log

import (
	"time"

	"github.com/form3tech/go-context/tracectx"

	"github.com/sirupsen/logrus"
)

type transactionTimeFormatter struct {
	delegate logrus.Formatter
}

func (t *transactionTimeFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if tt, ok := entry.Data[tracectx.TransactionTime.String()]; ok {
		if transactionTime, ok := tt.(time.Time); ok {
			// build a new entry avoiding mutating the data map, which may be used by concurrent loggers
			// ensuring to re-use other fields such as Time and Buffer, which were initialised by log before calling
			// the formatter.
			data := make(logrus.Fields, len(entry.Data)+1)
			for k, v := range entry.Data {
				data[k] = v
			}
			data[tracectx.TransactionDuration.String()] = time.Since(transactionTime).Milliseconds()
			entry = &logrus.Entry{
				Logger:  entry.Logger,
				Data:    data,
				Time:    entry.Time,
				Level:   entry.Level,
				Caller:  entry.Caller,
				Message: entry.Message,
				Buffer:  entry.Buffer,
				Context: entry.Context,
			}
		}
	}
	return t.delegate.Format(entry)
}

func NewTransactionTimeFormatter(delegate logrus.Formatter) logrus.Formatter {
	return &transactionTimeFormatter{
		delegate: delegate,
	}
}

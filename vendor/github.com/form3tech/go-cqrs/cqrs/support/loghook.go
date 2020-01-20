package support

import "github.com/sirupsen/logrus"

type LogHook struct {
	logs []*logrus.Entry
}

func NewLogHook() *LogHook {
	return &LogHook{}
}

func (l *LogHook) Get() []*logrus.Entry {
	return l.logs
}

func (l *LogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (l *LogHook) Fire(entry *logrus.Entry) error {
	l.logs = append(l.logs, entry)
	return nil
}

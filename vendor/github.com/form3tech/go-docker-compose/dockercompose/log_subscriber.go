package dockercompose

import log "github.com/sirupsen/logrus"

type LogSubscriber interface {
	OnNext(log string)
}

type logSubscriber struct {
	onNextFunc func(string)
}

func (s logSubscriber) OnNext(logLine string) {
	s.onNextFunc(logLine)
}

func ToLogSubscriber(onNextFunc func(logLine string)) LogSubscriber {
	return logSubscriber{
		onNextFunc: onNextFunc,
	}
}

var defaultLogSubscriber = ToLogSubscriber(func(logLine string) {
	log.Info(logLine)
})

func newLogSubscriber(subscribers ...LogSubscriber) LogSubscriber {
	subscribers = append(subscribers, defaultLogSubscriber)
	return ToLogSubscriber(func(logLine string) {
		for _, subscriber := range subscribers {
			subscriber.OnNext(logLine)
		}
	})
}

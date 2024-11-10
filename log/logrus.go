package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	*logrus.Logger
}

func (a *LogrusLogger) WithField(key string, value interface{}) Logger {
	return NewFromLogrusEntry(a.Logger.WithField(key, value))
}

func (a *LogrusLogger) WithFields(fields Fields) Logger {
	return NewFromLogrusEntry(a.Logger.WithFields(logrus.Fields(fields)))
}

func (a *LogrusLogger) SetLevel(level Level) {
	a.Logger.SetLevel(toLogrusLevel(level))
}

type LogrusEntry struct {
	*logrus.Entry
}

func (a *LogrusEntry) WithField(key string, value interface{}) Logger {
	return NewFromLogrusEntry(a.Entry.WithField(key, value))
}

func (a *LogrusEntry) WithFields(fields Fields) Logger {
	return NewFromLogrusEntry(a.Entry.WithFields(logrus.Fields(fields)))
}

func (a *LogrusEntry) SetLevel(level Level) {
	a.Entry.Logger.SetLevel(toLogrusLevel(level))
}

func (a *LogrusEntry) SetOutput(w io.Writer) {
	a.Entry.Logger.SetOutput(w)
}

func NewFromLogrus(logger *logrus.Logger) Logger {
	return &LogrusLogger{logger}
}

func NewFromLogrusEntry(e *logrus.Entry) Logger {
	return &LogrusEntry{e}
}

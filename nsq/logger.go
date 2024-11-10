package nsq

import (
	"context"
	"regexp"
	"strings"

	"github.com/niming-dev/ddd-demo/go-common/log"
)

type cunsumerLogger struct {
}

var re *regexp.Regexp = regexp.MustCompile(`^([0-9]+) \[([^/]+)/([^\]]+)\]`)

func (*cunsumerLogger) Output(calldepth int, s string) error {
	logFields := log.Fields{
		log.FieldKeyModule: "NSQ/C",
	}
	level := ""
	if len(s) > 2 {
		level = s[0:3]
		s = strings.TrimLeft(s[3:], " ")
	}

	if matched := re.FindStringSubmatch(s); matched != nil {
		logFields["instanceID"] = matched[1]
		logFields["topic"] = matched[2]
		logFields["channel"] = matched[3]
		s = strings.TrimLeft(s[len(matched[0]):], " ")
	}
	logger := log.WithFields(context.Background(), logFields)
	switch level {
	case "DBG":
		logger.Debug(s)
	case "INF":
		logger.Info(s)
	case "WRN":
		logger.Warn(s)
	case "ERR":
		logger.Error(s)
	default:
		logger.Info(nil, s)
	}
	return nil
}

package log

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// refer to https://github.com/sirupsen/logrus/blob/master/text_formatter.go

const (
	defaultTimestampFormat = time.RFC3339
)

type Coloring func(format string, a ...interface{}) string

var (
	blue   Coloring = color.BlueString
	gray   Coloring = color.HiBlackString
	green  Coloring = color.GreenString
	yellow Coloring = color.YellowString
	red    Coloring = color.RedString
	cyan   Coloring = color.CyanString
)

// TextFormatter formats logs into text
type TextFormatter struct {
	// TimestampFormat to use for display when a full timestamp is printed.
	// The format to use is the same than for time.Format or time.Parse from the standard
	// library.
	// The standard Library already provides a set of predefined format.
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// The keys sorting function, when uninitialized it uses sort.Strings.
	SortingFunc func([]string)

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)
}

// Format renders a single log entry
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(Fields)
	for k, v := range entry.Data {
		data[k] = v
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.appendFixedValue(b, entry.Time.Format(timestampFormat))
	f.appendFixedValue(b, coloredLevel(entry.Level))

	fixedKeys := []string{
		FieldKeyModule,
		FieldKeyTraceID,
		FieldKeySpanID,
		FieldKeyParentSpanID,
		FieldKeyElapsed,
	}
	for _, key := range fixedKeys {
		if value, ok := data[string(key)]; ok {
			delete(data, key)
			f.appendKeyValue(b, cyan(key), value)
		}
	}

	var funcVal, fileVal string
	if entry.HasCaller() {
		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		} else {
			funcVal = entry.Caller.Function
			fileVal = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		}

		if funcVal != "" {
			f.appendKeyValue(b, logrus.FieldKeyFunc, funcVal)
		}
		if fileVal != "" {
			f.appendKeyValue(b, logrus.FieldKeyFile, fileVal)
		}
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	if !f.DisableSorting {
		if f.SortingFunc == nil {
			sort.Strings(keys)
		} else {
			f.SortingFunc(keys)
		}
	}

	for _, key := range keys {
		var value interface{}
		switch {
		case key == logrus.FieldKeyLogrusError:
			// FIXME unable to get entry.err
			// value = entry.err
		default:
			value = data[key]
		}
		f.appendKeyValue(b, key, value)
	}

	if entry.Message != "" {
		f.appendFixedValue(b, entry.Message)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(blue(key))
	b.WriteByte('=')
	f.appendValue(b, value, true)
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}, needsQuoting bool) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if needsQuoting {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	} else {
		b.WriteString(stringVal)
	}
}

func (f *TextFormatter) appendFixedValue(b *bytes.Buffer, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	f.appendValue(b, value, false)
}

func coloredLevel(level logrus.Level) string {
	var coloring Coloring
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		coloring = gray
	case logrus.WarnLevel:
		coloring = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		coloring = red
	case logrus.InfoLevel:
		coloring = green
	default:
		coloring = blue
	}

	return coloring("%-7s", strings.ToUpper(level.String()))
}

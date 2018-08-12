package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func NewStdLogger() *StdLogger {
	return &StdLogger{
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

type StdLogger struct {
	logger *log.Logger

	fields []Field
	prefix string // Computed from fields
}

func (l *StdLogger) Errorf(format string, v ...interface{}) {
	l.logger.Printf("%s: %s", l.prefix, fmt.Sprintf(format, v...))
}

func (l *StdLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf("%s: %s", l.prefix, fmt.Sprintf(format, v...))
}

func (l *StdLogger) Println(v ...interface{}) {
	l.logger.Printf("%s: %s", l.prefix, fmt.Sprint(v...))
}

func (l *StdLogger) With(fields ...Field) Logger {
	newFields := append(l.fields, fields...)

	mapped := map[string]interface{}{}
	for i := range newFields {
		mapped[newFields[i].Key] = newFields[i].Value
	}

	prefix := l.prefix

	data, err := json.Marshal(mapped)
	if err != nil {
		prefix = prefix + "| Error when adding more fields |"
	} else {
		prefix = string(data)
	}

	return &StdLogger{
		logger: l.logger,
		fields: newFields,
		prefix: prefix,
	}
}

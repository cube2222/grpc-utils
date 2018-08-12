package logger

import (
	"context"
	"fmt"
)

type Logger interface {
	Errorf(format string, v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	With(...Field) Logger
}

type Field struct {
	Key   string
	Value interface{}
}

func NewField(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

func (f Field) String() string {
	return fmt.Sprintf("%s: %v", f.Key, f.Value)
}

func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value("logger").(Logger)
	if !ok || logger == nil {
		return NewStdLogger()
	}

	return logger
}

func Inject(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

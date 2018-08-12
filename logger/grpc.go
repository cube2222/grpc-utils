package logger

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func GRPCInjector(log Logger, keys ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var fields []Field

		for _, key := range keys {
			value := ctx.Value(key)
			fields = append(fields, NewField(key, value))
		}

		curRequestLogger := log.With(fields...)
		ctx = Inject(ctx, curRequestLogger)

		return handler(ctx, req)
	}
}

func GRPCServerLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		curRequestLogger := FromContext(ctx)

		start := time.Now()

		res, err := handler(ctx, req)

		duration := time.Since(start)

		if err != nil {
			s, ok := status.FromError(err)
			if ok {
				curRequestLogger = curRequestLogger.With(
					NewField("err.msg", s.Message()),
					NewField("err.code", s.Code()),
				)
			} else {
				curRequestLogger = curRequestLogger.With(
					NewField("err", err.Error()),
				)
			}
		}

		curRequestLogger.With(
			NewField("path", info.FullMethod),
			NewField("duration", duration),
		).Printf("Finished handling request.")

		return res, err
	}
}

func GRPCClientLogger() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		curRequestLogger := FromContext(ctx)

		start := time.Now()

		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(start)

		if err != nil {
			s, ok := status.FromError(err)
			if ok {
				curRequestLogger = curRequestLogger.With(
					NewField("err.msg", s.Message()),
					NewField("err.code", s.Code()),
				)
			} else {
				curRequestLogger = curRequestLogger.With(
					NewField("err", err.Error()),
				)
			}
		}

		curRequestLogger.With(
			NewField("path", method),
			NewField("duration", duration),
		).Printf("Finished doing request.")

		return err
	}
}

package requestid

import (
	"context"
	"log"
	"net/http"

	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const Key = "request-id"

func ServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var requestID string

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			value := md.Get(Key)
			if len(value) == 1 {
				requestID = value[0]
			}
		}

		if requestID == "" {
			requestID = GenerateRequestID()
		}

		ctx = context.WithValue(ctx, Key, requestID)
		return handler(ctx, req)
	}
}

func ClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var requestID string
		if value, ok := ctx.Value(Key).(string); ok {
			requestID = value
		}

		if requestID == "" {
			requestID = GenerateRequestID()
		}

		ctx = metadata.AppendToOutgoingContext(ctx, Key, requestID)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func HTTPInterceptor(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(Key)
		if requestID == "" {
			requestID = GenerateRequestID()
		}

		r = r.WithContext(context.WithValue(r.Context(), Key, requestID))
		h.ServeHTTP(w, r)
	})
}

func GenerateRequestID() string {
	var requestID string

	id, err := uuid.NewV4()
	if err != nil {
		log.Printf("Couldn't generate request-id: %v", err)
		requestID = "err"
	}
	requestID = id.String()

	return requestID
}

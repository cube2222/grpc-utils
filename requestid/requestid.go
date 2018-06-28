package requestid

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			t := grpc_ctxtags.Extract(ctx)
			if !t.Has("request-id") {
				t.Set("request-id", md["request-id"][0])
			} else {
				log.Println("Missing request-id")
			}
		} else {
			log.Println("Missing metadata")
		}
		return handler(ctx, req)
	}
}

func ClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ID, ok := grpc_ctxtags.Extract(ctx).Values()["request-id"].(string)
		if !ok || ID == "" {
			// Fall back to taking out of context
			ID, ok = ctx.Value("request-id").(string)
			if !ok {
				newID, err := uuid.NewV4()
				if err != nil {
					log.Println("Couldn't generate request uuid: ", err)
				}
				ID = newID.String()
			}
		}

		md := metadata.Pairs("request-id", ID)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func HTTPInjector(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID, err := uuid.NewV4()
		if err != nil {
			log.Println("Couldn't generate request uuid: ", err)
		}
		r = r.WithContext(context.WithValue(r.Context(), "request-id", ID.String()))
		h.ServeHTTP(w, r)
	})
}

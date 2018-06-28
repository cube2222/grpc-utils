package httplogger

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

func HTTPInject(h http.Handler) http.Handler {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Couldn't create http request zap logger.")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		curRequestLogger := logger

		requestID, ok := r.Context().Value("request-id").(string)
		if ok {
			curRequestLogger = logger.With(
				zap.String("request-id", requestID),
			)
		}
		r = r.WithContext(ctxzap.ToContext(r.Context(), curRequestLogger))

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		start := time.Now()
		h.ServeHTTP(ww, r)
		duration := time.Since(start)

		curRequestLogger.Info("Finished request.",
			zap.String("user-agent", r.UserAgent()),
			zap.String("host", r.Host),
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
			zap.Int("code", ww.Status()),
			zap.Int("bytes-written", ww.BytesWritten()),
		)
	})
}

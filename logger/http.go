package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

func HTTPInjector(log Logger, keys ...string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var fields []Field

			for _, key := range keys {
				value := r.Context().Value(key)
				fields = append(fields, NewField(key, value))
			}

			curRequestLogger := FromContext(r.Context()).With(fields...)

			r = r.WithContext(Inject(r.Context(), curRequestLogger))

			h.ServeHTTP(w, r)
		})
	}
}

func HTTPLogger() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			h.ServeHTTP(ww, r)
			duration := time.Since(start)

			FromContext(r.Context()).With(
				NewField("user-agent", r.UserAgent()),
				NewField("host", r.Host),
				NewField("path", r.URL.Path),
				NewField("method", r.Method),
				NewField("duration", duration),
				NewField("code", ww.Status()),
			).Printf("Finished handling request.")
		})
	}
}

package middleware

import (
	"gaman-microservice/api-gateway/constant"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type statusCapturingWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusCapturingWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusCapturingWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

func (w *statusCapturingWriter) IsError() bool {
	return w.status >= 400
}

func LoggerMiddleware() runtime.Middleware {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			now := time.Now()
			reqId, ok := r.Context().Value(constant.CtxKeyRequestID).(string)
			if !ok {
				reqId = "unknown"
			}

			// inject logger to context
			l := log.With().
				Str("request_id", reqId).
				Logger()
			r = r.WithContext(l.WithContext(r.Context()))

			l.Info().
				Any("pathParams", pathParams).
				Msgf("[http request] %s %s - %s", r.Method, r.URL.Path, r.UserAgent())

			sw := &statusCapturingWriter{ResponseWriter: w}

			next(sw, r, pathParams)

			level := zerolog.InfoLevel
			if sw.IsError() {
				level = zerolog.ErrorLevel
			}
			l.WithLevel(level).
				Msgf("[http response] %d %s", sw.Status(), time.Since(now))
		}
	}
}

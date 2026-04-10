package loggedHttpclient

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type loggingRoundTripper struct {
	transport http.RoundTripper
}

func (l *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	logger := log.Ctx(req.Context()).With().
		Str("method", req.Method).
		Str("url", req.URL.String()).
		Logger()

	logger.Info().Msg("sending HTTP request")

	start := time.Now()
	resp, err := l.transport.RoundTrip(req)
	duration := time.Since(start)

	if err != nil {
		logger.Error().
			Err(err).
			Dur("duration", duration).
			Msg("HTTP request failed")
		return resp, err
	}

	logger.
		Info().
		Int("status_code", resp.StatusCode).
		Dur("duration", duration).
		Msg("received HTTP response")

	return resp, nil
}

func New() *http.Client {
	return &http.Client{
		Transport: &loggingRoundTripper{
			transport: http.DefaultTransport,
		},
		Timeout: 30 * time.Second,
	}
}

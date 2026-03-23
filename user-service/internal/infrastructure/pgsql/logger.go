package pgsql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

const slowQueryThreshold = 2 * time.Second

type queryTracerKey struct{}

var tracerKey = queryTracerKey{}

type queryData struct {
	timeStart time.Time
	query     string
	args      []any
}

type logger struct {
}

func (l logger) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	qd := &queryData{
		timeStart: time.Now(),
		query:     data.SQL,
		args:      data.Args,
	}

	return context.WithValue(ctx, tracerKey, qd)
}

func (l logger) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	qd, ok := ctx.Value(tracerKey).(*queryData)
	if !ok {
		return
	}

	duration := time.Since(qd.timeStart)

	level := zerolog.InfoLevel
	if data.Err != nil {
		level = zerolog.ErrorLevel
	} else if duration > slowQueryThreshold {
		level = zerolog.WarnLevel
	}

	zerolog.Ctx(ctx).
		WithLevel(level).
		Err(data.Err).
		Str("duration", duration.String()).
		Any("args", qd.args).
		Msgf("Query: %s", qd.query)
}

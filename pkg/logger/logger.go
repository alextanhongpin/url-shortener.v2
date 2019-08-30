package logger

import (
	"context"

	"go.uber.org/zap"
)

func WithContext(ctx context.Context, fields ...zap.Field) *zap.Logger {
	// TODO: Populate with request id.
	// fields = append(fields, zap.String("req_id", ""))
	return zap.L().With(fields...)
}

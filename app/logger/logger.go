package logger

import (
	"context"

	"github.com/alextanhongpin/pkg/requestid"

	"go.uber.org/zap"
)

func WithContext(ctx context.Context, fields ...zap.Field) *zap.Logger {
	reqid, ok := requestid.Value(ctx)
	if ok {
		fields = append(fields, zap.String("req_id", reqid))
	}
	return zap.L().With(fields...)
}

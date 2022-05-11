package logger

import (
	"context"

	"gin-boilerplate/comm/trace"
)

func Extract(ctx context.Context) (traceID string, parentSpanID string, isFound bool) {
	return trace.FromContext(ctx)
}

func ExtractTraceID(ctx context.Context) string {
	traceID, _, _ := Extract(ctx)
	return traceID
}

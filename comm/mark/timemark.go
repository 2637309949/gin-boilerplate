package mark

import (
	"context"
	"fmt"
	"gin-boilerplate/comm/logger"
	"time"
)

type TimeMark struct {
	ctx         context.Context
	start, last time.Time
}

func (t *TimeMark) Mark(format string, data ...interface{}) {
	duration := time.Since(t.last)
	desc := fmt.Sprintf(format, data...)
	logger.Infof(t.ctx, "%s duration:[%s]", desc, duration)
	t.last = time.Now()
}

func (t *TimeMark) Init(ctx context.Context, format string, data ...interface{}) func() {
	t.ctx = ctx
	t.start = time.Now()
	t.last = time.Now()
	return func() {
		desc := fmt.Sprintf(format, data...)
		duration := time.Since(t.start)
		logger.Infof(t.ctx, "%s total duration:[%v]", desc, duration)
	}
}

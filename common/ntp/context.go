package ntp

import (
	"context"
	"time"

	"github.com/khayyamov/sing/service"
)

func TimeFuncFromContext(ctx context.Context) func() time.Time {
	timeService := service.FromContext[TimeService](ctx)
	if timeService == nil {
		return nil
	}
	return timeService.TimeFunc()
}

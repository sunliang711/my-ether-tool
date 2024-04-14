package utils

import (
	"context"
	"met/consts"
	"time"
)

func DefaultTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*consts.DefaultTimeout)
}

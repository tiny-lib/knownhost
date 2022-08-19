package knownhost

import (
	"context"
	"time"
)

const (
	hostKey      = "host"
	algorithmKey = "algorithm"
	timeoutKey   = "timeout"
)

func getHostFromContext(ctx context.Context) string {
	return ctx.Value(hostKey).(string)
}

func getAlgorithmFromContext(ctx context.Context) string {
	return ctx.Value(algorithmKey).(string)
}

func getTimeoutFromContext(ctx context.Context) time.Duration {
	return ctx.Value(timeoutKey).(time.Duration)
}

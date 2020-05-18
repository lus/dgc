package dgc

import (
	"time"

	"github.com/zekroTJA/timedmap"
)

// RateLimiter represents a rate limiter
type RateLimiter struct {
	Cooldown           time.Duration
	RateLimitedHandler ExecutionHandler
	executions         *timedmap.TimedMap
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(cooldown, cleanupInterval time.Duration, onRateLimited ExecutionHandler) *RateLimiter {
	return &RateLimiter{
		Cooldown:           cooldown,
		RateLimitedHandler: onRateLimited,
		executions:         timedmap.New(cleanupInterval),
	}
}

// NotifyExecution notifies the rate limiter about a new execution and returns whether or not the execution is allowed
func (rateLimiter *RateLimiter) NotifyExecution(ctx *Ctx) bool {
	if rateLimiter.executions.Contains(ctx.Event.Author.ID) {
		if rateLimiter.RateLimitedHandler != nil {
			nextExecution, err := rateLimiter.executions.GetExpires(ctx.Event.Author.ID)
			if err != nil {
				ctx.CustomObjects["dgc_nextExecution"] = nextExecution
			}
			rateLimiter.RateLimitedHandler(ctx)
		}
		return false
	}
	rateLimiter.executions.Set(ctx.Event.Author.ID, time.Now().UnixNano()/1e6, rateLimiter.Cooldown, nil)
	return true
}

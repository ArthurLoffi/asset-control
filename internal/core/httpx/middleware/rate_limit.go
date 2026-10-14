package middleware

import (
	"asset-control/internal/core/exception"
	"asset-control/internal/core/httpx"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimit allows each client IP a burst of 15 requests, refilled at 2 per
// second. Every call owns its visitors, so separate routers never share budgets;
// idle visitors are swept on access instead of by a background goroutine.
func RateLimit() gin.HandlerFunc {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)
	lastSweep := time.Now()
	return func(c *gin.Context) {
		now := time.Now()
		ip := c.ClientIP()
		mu.Lock()
		if now.Sub(lastSweep) > time.Minute {
			for address, idle := range visitors {
				if now.Sub(idle.lastSeen) > 3*time.Minute {
					delete(visitors, address)
				}
			}
			lastSweep = now
		}
		current, exists := visitors[ip]
		if !exists {
			current = &visitor{limiter: rate.NewLimiter(2, 15)}
			visitors[ip] = current
		}
		current.lastSeen = now
		isAllowed := current.limiter.Allow()
		mu.Unlock()
		if !isAllowed {
			httpx.WriteError(c, exception.NewHTTPError(http.StatusTooManyRequests, "Too many requests. Try again later."))
			return
		}
		c.Next()
	}
}

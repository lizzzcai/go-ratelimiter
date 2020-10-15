package ratelimiter

import (
	"sync/atomic"
	"time"
)

type Counter struct {
	Count       uint64 // init the counter
	Limit       uint64 // max number of requests in window
	Interval    int64  // unit: ms
	RefreshTime int64  // time window
}

func NewCounter(count, limit uint64, interval, rt int64) *Counter {
	return &Counter{
		Count:       count,
		Limit:       limit,
		Interval:    interval,
		RefreshTime: rt,
	}
}

func (c *Counter) RateLimit() bool {
	now := time.Now().UnixNano() / 1e6
	if now < (c.RefreshTime + c.Interval) {
		atomic.AddUint64(&c.Count, 1)
		return c.Count <= c.Limit
	} else {
		c.RefreshTime = now
		atomic.AddUint64(&c.Count, -c.Count)
		return true
	}
}

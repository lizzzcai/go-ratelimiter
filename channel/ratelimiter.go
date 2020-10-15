package main

import (
	"errors"
	"fmt"
	"time"
)

type RateLimiter struct {
	limit    float64 // number of tokens generated per second
	burst    int     // bucket size
	interval time.Duration
	bucket   chan time.Time
	tokens   <-chan time.Time
}

func NewRateLimiter(limit float64, burst int) (*RateLimiter, error) {
	if limit <= 0 || burst <= 0 {
		return nil, errors.New("Limit or burst can not be negative.")
	}
	bucket := make(chan time.Time, burst)
	for i := 0; i < burst; i++ {
		bucket <- time.Now()
	}

	duration := time.Duration(1000 / limit)
	tokens := time.Tick(duration * time.Millisecond)

	return &RateLimiter{
		limit,
		burst,
		duration,
		bucket,
		tokens,
	}, nil
}

func (rl *RateLimiter) Start() {
	go func() {
		for t := range rl.tokens {
			select {
			case rl.bucket <- t:
			default:
			}
		}
	}()
}

func (rl *RateLimiter) Acquire() error {
	select {
	case <-rl.bucket:
		return nil
	case <-time.After(rl.interval * time.Millisecond):
		return errors.New("time out")
	default:
		return errors.New("no token")
	}
}

func main() {

	ratelimiter, err := NewRateLimiter(5, 3)
	if err != nil {
		fmt.Println(err)
		return
	}
	ratelimiter.Start()

	i := 0
	for range time.Tick(100 * time.Millisecond) {
		if err := ratelimiter.Acquire(); err == nil {
			fmt.Println("request", i, time.Now())
		} else {
			fmt.Println("request", i, err)
		}
		i += 1
		if i > 20 {
			break
		}
	}

}

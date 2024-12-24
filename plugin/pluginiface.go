package plugin

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type PluginInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type leakyBucketRateLimit struct {
	MaxRequests int
	Duration    int
	userBuckets map[string]*userBucket
}

type userBucket struct {
	count int
	mutex sync.Mutex
}

func NewLeakyBucketRateLimit(maxRequests int, duration int) PluginInterface {

	limiter := &leakyBucketRateLimit{
		MaxRequests: maxRequests,
		Duration:    duration,
		userBuckets: map[string]*userBucket{},
	}
	timer := time.Tick(time.Duration(duration) * time.Second)

	go func(timer <-chan time.Time) {
		for range timer {

			fmt.Println("Refilling buckets")
			for user, bucket := range limiter.userBuckets {

				bucket.mutex.Lock()
				if bucket.count < limiter.MaxRequests {
					fmt.Println("Refilling buckets of user", user)
					bucket.count++
				}
				bucket.mutex.Unlock()
			}
		}
	}(timer)

	return limiter

}

func (l *leakyBucketRateLimit) Handle(w http.ResponseWriter, r *http.Request) {
	// TODO implement rate limit

	userId := r.Header["X-User-Id"]

	if !l.isAllowed(userId[0]) {
		http.Error(w, "too Many Request", http.StatusTooManyRequests)
		return
	}
}

func (l leakyBucketRateLimit) isAllowed(user string) bool {

	buckets, found := l.userBuckets[user]
	if !found {

		l.userBuckets[user] = &userBucket{
			count: l.MaxRequests - 1,
			mutex: sync.Mutex{},
		}
		slog.Debug("created new bucket for user", slog.String("userId", user))
		return true
	}

	slog.Debug("State:", slog.String("user", user), slog.Int("buckets", buckets.count))

	if buckets.count > 0 {
		buckets.mutex.Lock()
		buckets.count--
		buckets.mutex.Unlock()
	}

	return buckets.count > 0

}

package plugin

import (
	"fmt"
	"gateor/pkg"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type leakyBucketRateLimit struct {
	MaxRequests int
	Duration    int
	userBuckets map[string]*userBucket
	NextPlugin  PluginInterface
	mu          sync.Mutex
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

	go limiter.refill(time.Tick(10 * time.Second))

	return limiter

}

func (l *leakyBucketRateLimit) Handle(w http.ResponseWriter, r *http.Request) {

	apiKey, ok := r.Header["X-Api-Key"]
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !l.isAllowed(apiKey[0]) {
		http.Error(w, "too Many Request", http.StatusTooManyRequests)
		return
	}
	if l.NextPlugin != nil {
		l.NextPlugin.Handle(w, r)
	}

}

func (l *leakyBucketRateLimit) AddNext(nextPlugin PluginInterface) {
	l.NextPlugin = nextPlugin
}

func (limiter *leakyBucketRateLimit) refill(timer <-chan time.Time) {

	for t := range timer {
		pkg.Log.Debug("TICK", slog.Any("time", t))
		limiter.mu.Lock()
		for user, bucket := range limiter.userBuckets {

			bucket.mutex.Lock()
			if bucket.count < limiter.MaxRequests {
				fmt.Println("Refilling buckets of user", user)
				bucket.count++
			}
			bucket.mutex.Unlock()
		}
		limiter.mu.Unlock()
	}
}

func (l *leakyBucketRateLimit) isAllowed(user string) bool {

	l.mu.Lock()
	defer l.mu.Unlock()

	buckets, found := l.userBuckets[user]

	// New client, create new bucket
	if !found {
		l.userBuckets[user] = &userBucket{
			count: l.MaxRequests - 1,
			mutex: sync.Mutex{},
		}
		slog.Debug("created new bucket for user", slog.String("userId", user))
		return true
	}

	// Existing client
	pkg.Log.Debug("State:", slog.String("user", user), slog.Int("buckets", buckets.count))

	// client have more than 0 requests
	if buckets.count > 0 {
		buckets.mutex.Lock()
		buckets.count--
		buckets.mutex.Unlock()
		return true
	}

	// client has 0 requests left
	return false

}

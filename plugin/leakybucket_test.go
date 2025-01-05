package plugin

import (
	"fmt"
	"sync"
	"testing"
)

func TestLeakyBucketRateLimit_IsAllowed_Concurrent(t *testing.T) {
	// Initialize the rate limiter
	limiter := &leakyBucketRateLimit{
		MaxRequests: 10,
		Duration:    60,
		userBuckets: map[string]*userBucket{},
	}

	// Number of goroutines
	numGoroutines := 100

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Function to be executed by each goroutine
	testFunc := func(user string) {
		defer wg.Done()
		for i := 0; i < 15; i++ {
			if !limiter.isAllowed(user) {
				t.Errorf("Request denied for user %s", user)
			}
		}
	}

	// Start multiple goroutines
	for i := 0; i < numGoroutines; i++ {
		go testFunc(fmt.Sprintf("user-%d", i))
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

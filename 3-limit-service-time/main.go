//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"context"
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	mu        sync.Mutex
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequestBeginner(process func(), u *User) bool {
	if u.IsPremium == false {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
		defer cancel()

		processDone := make(chan struct{})

		go func() {
			process()
			close(processDone)
		}()

		select {
		case <-processDone:
			return true
		case <-ctx.Done():
			return false
		}
	}

	process()
	return true
}

func HandleRequest(process func(), u *User) bool {
	if u.IsPremium == false {
		u.mu.Lock()
		defer u.mu.Unlock()

		remaining := 10 - u.TimeUsed
		if remaining <= 0 {
			return false
		}
		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(time.Second*time.Duration(remaining)),
		)

		defer cancel()

		processDone := make(chan struct{})

		go func() {
			start := time.Now()
			process()
			elapsed := time.Since(start)
			u.TimeUsed += int64(elapsed.Seconds())
			close(processDone)
		}()

		select {
		case <-processDone:
			return true
		case <-ctx.Done():
			return false
		}
	}

	process()
	return true
}

func main() {
	RunMockServer()
}

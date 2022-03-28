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
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	mux       sync.Mutex
}

const maxFreeTime = 10

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	// if !u.IsPremium && u.TimeUsed > 10 {
	// 	return false
	// }
	// start := time.Now()
	// process()
	// elapsed := time.Now().Sub(start).Seconds()
	// u.TimeUsed += int64(elapsed)
	// return u.IsPremium || u.TimeUsed <= 10

	done := make(chan bool)
	go func(chan<- bool) {
		process()
		done <- true
	}(done)
	for {
		select {
		case <-done:
			return true
		case <-time.Tick(time.Second):
			{
				u.mux.Lock()

				u.TimeUsed++

				if !u.IsPremium && u.TimeUsed > maxFreeTime {
					u.mux.Unlock()
					return false
				}
				u.mux.Unlock()

			}

		}
	}

}

func main() {
	RunMockServer()
}

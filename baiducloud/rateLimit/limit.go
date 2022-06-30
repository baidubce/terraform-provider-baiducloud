package ratelimit

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var (
	locker              sync.Mutex
	bccPurchaseLimitNum = int32(0)
)

func Check(action string) error {

	old := time.Now()
	for true {
		locker.Lock()
		if bccPurchaseLimitNum <= DefaultLimit {
			locker.Unlock()
			break
		}
		locker.Unlock()
		intn := rand.Intn(200) + 200
		time.Sleep(time.Millisecond * time.Duration(intn))
		if time.Since(old) > 5*time.Minute {
			return fmt.Errorf("%s wait too long :exceed 5 minute ", action)
		}
	}
	atomic.AddInt32(&bccPurchaseLimitNum, 1)
	return nil
}

func CheckEnd() {
	atomic.AddInt32(&bccPurchaseLimitNum, -1)
}

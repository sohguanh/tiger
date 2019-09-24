// rateLimiter is the package that contain the different algorithm for API rate limiting.
package rateLimiter

import (
	"io"
	"net/http"
	"strconv"
	"sync"
	logUtil "tiger/util/log"
	"time"
)

// TokenBucketHandler is the struct that contain all information needed for Token Bucket algorithm.
type TokenBucketHandler struct {
	Lock              sync.Mutex
	MaximumAmt        int
	CurrentAmt        int
	RefillAmt         int
	RefillSpawn       bool
	Ticker            *time.Ticker
	TickerDurationSec int
	TickerQuitChan    chan bool
}

// NewTokenBucketHandler is to create a new Token Bucket handler. maximumAmt parameter to indicate how many request the url endpoint can handle before rejecting. refillDurationSec parameter to indicate elapsed how many seconds before refilling. refillAmt to indicate how many to refill. any value that exceed maximumAmt will still be capped to maximumAmt.
func NewTokenBucketHandler(maximumAmt int, refillDurationSec int, refillAmt int) *TokenBucketHandler {
	return &TokenBucketHandler{
		Lock:              sync.Mutex{},
		MaximumAmt:        maximumAmt,
		CurrentAmt:        maximumAmt,
		RefillAmt:         refillAmt,
		RefillSpawn:       false,
		Ticker:            nil,
		TickerDurationSec: refillDurationSec,
		TickerQuitChan:    make(chan bool),
	}
}

// ServeNextHTTP is implementation method for the ChainNextHandler interface.
func (a *TokenBucketHandler) ServeNextHTTP(w http.ResponseWriter, r *http.Request) bool {
	a.Lock.Lock()
	defer a.Lock.Unlock()
	if !a.RefillSpawn {
		a.RefillSpawn = true
		go func(handler *TokenBucketHandler) {
			logUtil.DebugPrintln("enter time ticker")
			a.Ticker = time.NewTicker(time.Duration(a.TickerDurationSec) * time.Second)
		LOOP:
			for {
				select {
				case <-a.Ticker.C:
					a.Lock.Lock()
					if a.CurrentAmt < a.MaximumAmt {
						newAmt := a.CurrentAmt + a.RefillAmt
						if newAmt > a.MaximumAmt {
							a.CurrentAmt = a.MaximumAmt
						} else {
							a.CurrentAmt = newAmt
						}
					}
					logUtil.DebugPrintln("CurrentAmt: " + strconv.Itoa(a.CurrentAmt))
					a.Lock.Unlock()
				case <-a.TickerQuitChan:
					break LOOP
				}
			}
			logUtil.DebugPrintln("exit time ticker")
		}(a)
	}

	if a.CurrentAmt > 0 {
		a.CurrentAmt = a.CurrentAmt - 1
		io.WriteString(w, "TokenBucketHandler handler accepted\n")
		io.WriteString(w, "CurrentAmt: "+strconv.Itoa(a.CurrentAmt)+"\n")
		//logUtil.DebugPrintln("CurrentAmt: "+strconv.Itoa(a.CurrentAmt))
		//logUtil.DebugPrintln("TokenBucketHandler handler accepted")
		return true
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "TokenBucketHandler handler rejected\n")
		io.WriteString(w, "CurrentAmt: "+strconv.Itoa(a.CurrentAmt)+"\n")
		//logUtil.DebugPrintln("CurrentAmt: "+strconv.Itoa(a.CurrentAmt))
		//logUtil.DebugPrintln("TokenBucketHandler handler rejected")
		return false
	}
}

// StopTimer is to stop the periodic refill timer from continuing. Not calling this will result in the timer that run "forever" periodically until server shut down.
func (a *TokenBucketHandler) StopTimer() {
	a.Lock.Lock()
	defer a.Lock.Unlock()
	if a.Ticker != nil {
		a.TickerQuitChan <- true
		a.Ticker.Stop()
		a.Ticker = nil
	}
}

package rateLimiter

import (
	"container/list"
	"io"
	"net/http"
	"sync"
	logUtil "tiger/util/log"
	"time"
)

// SlidingWindowHandler is the struct that contain all information needed for Sliding Window algorithm.
type SlidingWindowHandler struct {
	Lock          sync.Mutex
	RequestList   *list.List //at most 2 Element and within 1 minute interval
	RequestPerMin int
}

type slidingWindowCounter struct {
	count      int
	minuteMark int
}

// NewSlidingWindowHandler is to create a new Sliding Window handler. requestPerMin parameter to indicate how many request per minute the url endpoint can handle before rejecting.
func NewSlidingWindowHandler(requestPerMin int) *SlidingWindowHandler {
	return &SlidingWindowHandler{
		Lock:          sync.Mutex{},
		RequestList:   list.New(),
		RequestPerMin: requestPerMin,
	}
}

// ServeNextHTTP is implementation method for the ChainNextHandler interface.
func (a *SlidingWindowHandler) ServeNextHTTP(w http.ResponseWriter, r *http.Request) bool {
	if a.allowed() {
		io.WriteString(w, "SlidingWindowHandler handler accepted\n")
		//logUtil.DebugPrintln("SlidingWindowHandler handler accepted")
		return true
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "SlidingWindowHandler handler rejected\n")
		//logUtil.DebugPrintln("SlidingWindowHandler handler rejected")
		return false
	}
}

func (a *SlidingWindowHandler) allowed() bool {
	a.Lock.Lock()
	defer a.Lock.Unlock()

	_, currentMin, currentSec := time.Now().Clock()
	if a.RequestList.Len() == 0 {
		logUtil.DebugPrintln("a.RequestList.Len() == 0")
		a.RequestList.PushBack(&slidingWindowCounter{count: 1, minuteMark: currentMin})
		return true
	} else if a.RequestList.Len() <= 2 {
		logUtil.DebugPrintf("a.RequestList.Len() == %d\n", a.RequestList.Len())
		var front = a.RequestList.Front().Value.(*slidingWindowCounter)

		if front.minuteMark == currentMin { //within same bucket
			logUtil.DebugPrintln("within same bucket")
			if a.RequestList.Len() == 2 { //check previous bucket apply formula
				var back = a.RequestList.Back().Value.(*slidingWindowCounter)
				if weightFormula(back.count, float32(currentSec/60), front.count) > a.RequestPerMin {
					logUtil.DebugPrintln("within same bucket reject 0")
					return false
				}
			}
			newCount := front.count + 1
			if newCount > a.RequestPerMin {
				logUtil.DebugPrintf("within same bucket reject 1 with newCount %d\n", newCount)
				return false
			} else {
				front.count = newCount
				return true
			}
		} else { //not found in current bucket
			logUtil.DebugPrintln("not same bucket")
			var prevMinuteMark = currentMin - 1
			if currentMin == 0 {
				prevMinuteMark = 59
			}
			if front.minuteMark == prevMinuteMark {
				//check if allowed using weightFormula
				if weightFormula(front.count, float32(currentSec/60), 0) > a.RequestPerMin {
					logUtil.DebugPrintln("not same bucket reject 0")
					return false
				} else {
					if a.RequestList.Len() == 1 {
						//if Len() == 1 add new bucket in front
						a.RequestList.PushFront(&slidingWindowCounter{count: 1, minuteMark: currentMin})
					} else if a.RequestList.Len() == 2 {
						//if Len() == 2
						//copy front bucket to back bucket
						//overwrite front bucket with new counter value
						back := a.RequestList.Back().Value.(*slidingWindowCounter)

						back.count = front.count
						back.minuteMark = front.minuteMark

						front.count = 1
						front.minuteMark = currentMin
					}
					return true
				}
			} else {
				if a.RequestList.Len() == 1 {
					// if Len()==1 > 1 min interval so just overwrite front bucket with new counter value
					front.count = 1
					front.minuteMark = currentMin
				} else if a.RequestList.Len() == 2 {
					// if Len()==2 > 1 min interval so just overwrite front bucket with new counter value
					//remove the back bucket as not needed for compare anymore
					front.count = 1
					front.minuteMark = currentMin

					a.RequestList.Remove(a.RequestList.Back())
				}
				return true
			}
		}
	}
	logUtil.DebugPrintln("outside reject 999")
	return false //not supposed to come here but if it does return false by default
}

func weightFormula(prevCounter int, currMinuteMarkPct float32, currCounter int) int {
	return int(float32(prevCounter)*(1.-currMinuteMarkPct)) + currCounter
}

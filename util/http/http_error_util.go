package httpUtil

import (
	"log"
	"net/http"
	"os"
	"sync"
	logUtil "tiger/util/log"
)

type customHttpError struct {
	ErrorCode  int
	Error      string
	ErrorPage  string
	RespHeader map[string]string
}

var onceCustomHttpError sync.Once
var mapCustomHttpError map[int]*customHttpError
var mutexHttpError sync.RWMutex

// NotFound to show custom not found page if configured else revert to Go default.
func NotFound(w http.ResponseWriter, r *http.Request) {
	initHttpError()
	mutexHttpError.RLock()
	defer mutexHttpError.RUnlock()
	if value, found := mapCustomHttpError[http.StatusNotFound]; found {
		logUtil.DebugPrint("call error page " + value.ErrorPage)
		respHeader := value.RespHeader
		if respHeader != nil {
			for key, value := range respHeader {
				w.Header().Set(key, value)
			}
		}
		http.ServeFile(w, r, value.ErrorPage)
	} else {
		http.NotFound(w, r)
	}
}

// Error to show custom error page if configured else revert to Go default.
func Error(w http.ResponseWriter, r *http.Request, error string, code int) {
	initHttpError()
	mutexHttpError.RLock()
	defer mutexHttpError.RUnlock()
	if value, found := mapCustomHttpError[code]; found {
		logUtil.DebugPrint("call error page " + value.ErrorPage)
		respHeader := value.RespHeader
		if respHeader != nil {
			for key, value := range respHeader {
				w.Header().Set(key, value)
			}
		}
		http.ServeFile(w, r, value.ErrorPage)
	} else {
		http.Error(w, error, code)
	}
}

func initHttpError() {
	onceCustomHttpError.Do(func() { //singleton
		mapCustomHttpError = make(map[int]*customHttpError)
	})
}

// AddCustomErrorPage to show any application defined custom error page if configured else revert to Go default.
func AddCustomErrorPage(errorCode int, errorPage string, respHeader map[string]string) {
	initHttpError()
	mutexHttpError.Lock()
	defer mutexHttpError.Unlock()
	if _, err := os.Stat(errorPage); os.IsNotExist(err) {
		log.Print("error find error page " + errorPage)
		return
	}
	logUtil.DebugPrint("process " + errorPage)
	mapCustomHttpError[errorCode] = &customHttpError{ErrorCode: errorCode, Error: http.StatusText(errorCode), ErrorPage: errorPage, RespHeader: respHeader}
}

package httpUtil

import (
	"net/http"
	"regexp"
)

// ChainPathTokenHandler is the interface for application to implement for chaining Path Param url feature. Framework will process and pass in the placeholder values inside the pathParam map.
// bool return true or false to indicate to call the next handler or not
type ChainPathTokenHandler interface {
	ServeNextHTTP(w http.ResponseWriter, r *http.Request, pathParam map[string]string) bool
}

// ChainNextHandler is the interface for application to implement for chaining url feature.
// bool return true or false to indicate to call the next handler or not
type ChainNextHandler interface {
	ServeNextHTTP(w http.ResponseWriter, r *http.Request) bool
}

// AddChainHandler to add url mapping to handler. Direct url syntax. handler []ChainNextHandler where first handler will be processed then the next etc until the last handler.
func AddChainHandler(urlMapping string, handler []ChainNextHandler, httpVerb ...string) {
	initMapHandler()
	mutexHttp.Lock()
	defer mutexHttp.Unlock()
	addChainHandlerInternal(urlMapping, handler, nil, nil, nil, httpVerb...)
}

// AddChainHandlerRegEx to add url mapping to handler. Regular expression url syntax. handler []ChainNextHandler where first handler will be processed then the next etc until the last handler.
func AddChainHandlerRegEx(urlMapping string, handler []ChainNextHandler, httpVerb ...string) {
	initMapHandler()
	mutexHttp.Lock()
	defer mutexHttp.Unlock()
	re := getHandlerRe(urlMapping)
	if re == nil {
		return
	}
	addChainHandlerInternal(urlMapping, handler, nil, nil, re, httpVerb...)
}

// AddChainHandlerPathParam to add url mapping to handler. Placeholder syntax supported are {} and :
// 	Example {id} or :id
// handler []ChainPathTokenHandler where first handler will be processed then the next etc until the last handler.
func AddChainHandlerPathParam(urlMapping string, pathTokenHandler []ChainPathTokenHandler, httpVerb ...string) {
	initMapHandler()
	mutexHttp.Lock()
	defer mutexHttp.Unlock()
	pathToken := splitBySlashToken(urlMapping)
	addChainHandlerInternal(urlMapping, nil, pathTokenHandler, pathToken, nil, httpVerb...)
}

func addChainHandlerInternal(urlMapping string, handler []ChainNextHandler, pathTokenHandler []ChainPathTokenHandler, pathToken []string, re *regexp.Regexp, httpVerb ...string) {
	if len(httpVerb) == 0 { //default to http.MethodGet
		if re == nil && pathTokenHandler == nil {
			mapHandler[urlMapping] = &httpVerbHandler{HttpVerb: []string{http.MethodGet}, ChainNextHandler: handler}
		} else if pathTokenHandler != nil {
			mapHandlerPathParam[urlMapping] = &httpVerbHandler{HttpVerb: []string{http.MethodGet}, ChainPathTokenHandler: pathTokenHandler, PathToken: pathToken}
		} else if re != nil {
			mapHandlerRegEx[urlMapping] = &httpVerbHandler{HttpVerb: []string{http.MethodGet}, ChainNextHandler: handler, RegEx: re}
		}
		return
	}

	var verbs []string
	for _, verb := range httpVerb {
		if validHttpVerb[verb] {
			verbs = append(verbs, verb)
		}
	}
	if len(verbs) != 0 {
		if re == nil && pathTokenHandler == nil {
			mapHandler[urlMapping] = &httpVerbHandler{HttpVerb: verbs, ChainNextHandler: handler}
		} else if pathTokenHandler != nil {
			mapHandlerPathParam[urlMapping] = &httpVerbHandler{HttpVerb: verbs, ChainPathTokenHandler: pathTokenHandler, PathToken: pathToken}
		} else if re != nil {
			mapHandlerRegEx[urlMapping] = &httpVerbHandler{HttpVerb: verbs, ChainNextHandler: handler, RegEx: re}
		}
	}
}

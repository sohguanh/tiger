// httpUtil is the main package containing all http server related features for application usage.
//
// 	http_error_util.go
// 	Above package is for application to register their own custom http error code webpages. Optional.
//
// 	http_rewrite.go
// 	Above package is for application to enable url rewrite on the same http server. Optional
//
// 	http_util.go
// 	http_chain_util.go
// 	Above packages are for application to register their url and handler either as a single or a chain of handlers. Mandatory.
//
// 	handler_util.go
// 	Above file is the ENTRY POINT called by tiger framework for all application to add in their own application specific code. Functions inside this file act as placeholder for application to add. The keyword ENTRY POINT will be stated explicitly in the function documentation so take note.
package httpUtil

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"tiger/config"
	"tiger/util/log"
)

// PathTokenHandler is the interface for application to implement for Path Param url feature. Framework will process and pass in the placeholder values inside the pathParam map.
type PathTokenHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, pathParam map[string]string)
}

type httpVerbHandler struct {
	HttpVerb    []string
	NextHandler http.Handler
	RegEx       *regexp.Regexp

	//below to handle path param
	PathToken        []string
	PathTokenHandler *PathTokenHandler

	//below to handle ChainHandler*
	ChainNextHandler      []ChainNextHandler
	ChainPathTokenHandler []ChainPathTokenHandler
}

func (a *httpVerbHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpVerbFound := httpVerbOk(r, a.HttpVerb)
	if httpVerbFound {
		if a.NextHandler != nil {
			a.NextHandler.ServeHTTP(w, r)
		} else if a.ChainNextHandler != nil {
			for _, value := range a.ChainNextHandler {
				if ok := value.ServeNextHTTP(w, r); !ok {
					break
				}
			}
		} else {
			NotFound(w, r)
		}
	} else {
		NotFound(w, r)
	}
}

var validHttpVerb = map[string]bool{
	http.MethodGet:     true,
	http.MethodHead:    true,
	http.MethodPost:    true,
	http.MethodPut:     true,
	http.MethodPatch:   true,
	http.MethodDelete:  true,
	http.MethodConnect: true,
	http.MethodOptions: true,
	http.MethodTrace:   true,
}

var mutexHttp sync.RWMutex
var mux *http.ServeMux
var onceHttp sync.Once
var onceHandler sync.Once
var mapHandler map[string]http.Handler
var mapHandlerRegEx map[string]http.Handler
var mapHandlerPathParam map[string]http.Handler
var pathParamRE *regexp.Regexp
var pathParamSlashRE *regexp.Regexp

// NewServeMux to get a singleton customized http.ServeMutex oject
func NewServeMux(c *config.Config, db *sql.DB) *http.ServeMux {
	onceHttp.Do(func() { //singleton
		logUtil.DebugPrint("serve mux first time init\n")
		initMapHandler()
		mux = http.NewServeMux()
		setupRootHandler(c, db, mux)
		setupStaticPath(c, db, mux)
		for key, value := range mapHandler {
			mux.Handle(key, value)
		}
	})
	return mux
}

func initMapHandler() {
	onceHandler.Do(func() { //singleton
		mapHandler = make(map[string]http.Handler)
		mapHandlerRegEx = make(map[string]http.Handler)
		mapHandlerPathParam = make(map[string]http.Handler)
		pathParamRE, _ = regexp.Compile(`{\s*(\w+)\s*}|:\s*(\w+)`)
		pathParamSlashRE, _ = regexp.Compile("/+")
	})
}

func httpVerbOk(r *http.Request, httpVerb []string) bool {
	found := false
	for _, verb := range httpVerb {
		if r.Method == verb {
			found = true
			break
		}
	}
	return found
}

func getHandlerRe(urlMapping string) *regexp.Regexp {
	var re *regexp.Regexp
	var err error
	if re, err = regexp.Compile("(?i)" + urlMapping); err != nil { //ignore case
		return nil
	}
	return re
}

func splitBySlashToken(urlMapping string) []string {
	urlMapping = pathParamSlashRE.ReplaceAllString(urlMapping, "/")
	var pathToken []string
	for _, value := range strings.Split(urlMapping, "/") {
		if value != "" {
			pathToken = append(pathToken, value)
		}
	}
	return pathToken
}

// AddHandler to add url mapping to handler. Direct url syntax.
func AddHandler(urlMapping string, handler http.Handler, httpVerb ...string) {
	initMapHandler()
	mutexHttp.Lock()
	defer mutexHttp.Unlock()
	addHandlerInternal(urlMapping, handler, nil, nil, nil, httpVerb...)
}

// AddHandlerRegEx to add url mapping to handler. Regular expression url syntax.
func AddHandlerRegEx(urlMapping string, handler http.Handler, httpVerb ...string) {
	initMapHandler()
	mutexHttp.Lock()
	defer mutexHttp.Unlock()
	re := getHandlerRe(urlMapping)
	if re == nil {
		return
	}
	addHandlerInternal(urlMapping, handler, nil, nil, re, httpVerb...)
}

// AddHandlerPathParam to add url mapping to handler. Placeholder syntax supported are {} and :
// 	Example {id} or :id
func AddHandlerPathParam(urlMapping string, pathTokenHandler PathTokenHandler, httpVerb ...string) {
	initMapHandler()
	mutexHttp.Lock()
	defer mutexHttp.Unlock()
	pathToken := splitBySlashToken(urlMapping)
	addHandlerInternal(urlMapping, nil, &pathTokenHandler, pathToken, nil, httpVerb...)
}

func addHandlerInternal(urlMapping string, handler http.Handler, pathTokenHandler *PathTokenHandler, pathToken []string, re *regexp.Regexp, httpVerb ...string) {
	if len(httpVerb) == 0 { //default to http.MethodGet
		if re == nil && pathTokenHandler == nil {
			mapHandler[urlMapping] = &httpVerbHandler{HttpVerb: []string{http.MethodGet}, NextHandler: handler}
		} else if pathTokenHandler != nil {
			mapHandlerPathParam[urlMapping] = &httpVerbHandler{HttpVerb: []string{http.MethodGet}, PathTokenHandler: pathTokenHandler, PathToken: pathToken}
		} else if re != nil {
			mapHandlerRegEx[urlMapping] = &httpVerbHandler{HttpVerb: []string{http.MethodGet}, NextHandler: handler, RegEx: re}
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
			mapHandler[urlMapping] = &httpVerbHandler{HttpVerb: verbs, NextHandler: handler}
		} else if pathTokenHandler != nil {
			mapHandlerPathParam[urlMapping] = &httpVerbHandler{HttpVerb: verbs, PathTokenHandler: pathTokenHandler, PathToken: pathToken}
		} else if re != nil {
			mapHandlerRegEx[urlMapping] = &httpVerbHandler{HttpVerb: verbs, NextHandler: handler, RegEx: re}
		}
	}
}

func handleUrlPathEx(c *config.Config, db *sql.DB, mux *http.ServeMux, w http.ResponseWriter, r *http.Request) error {
	for key, value := range mapHandler {
		if found, _ := filepath.Match(key, r.URL.Path); found {
			logUtil.DebugPrintln("call match path url " + key)
			value.ServeHTTP(w, r)
			return nil
		}
	}
	return errors.New("cannot find match path url " + r.URL.Path)
}

func handleUrlRegEx(c *config.Config, db *sql.DB, mux *http.ServeMux, w http.ResponseWriter, r *http.Request) error {
	for key, value := range mapHandlerRegEx {
		if handler, found := value.(*httpVerbHandler); found {
			if handler.RegEx.MatchString(r.URL.Path) {
				logUtil.DebugPrintln("call match regex url " + key)
				value.ServeHTTP(w, r)
				return nil
			}
		}
	}
	return errors.New("cannot find match regex url " + r.URL.Path)
}

func handleUrlPathParam(c *config.Config, db *sql.DB, mux *http.ServeMux, w http.ResponseWriter, r *http.Request) error {
	actualToken := splitBySlashToken(r.URL.Path)

	for key, value := range mapHandlerPathParam {
		if handler, found := value.(*httpVerbHandler); found {
			var found bool = false
			var pathParam = make(map[string]string)
			if len(handler.PathToken) == len(actualToken) {
				for index, value := range handler.PathToken {
					if matched := pathParamRE.MatchString(value); matched {
						value1 := pathParamRE.ReplaceAllString(value, "$1$2")
						pathParam[value1] = actualToken[index]
						found = true
					} else {
						if value != actualToken[index] {
							found = false
							break
						}
					}
				}
			}
			if found {
				logUtil.DebugPrintln("call match path param url " + key)
				httpVerbFound := httpVerbOk(r, handler.HttpVerb)
				if httpVerbFound {
					if handler.PathTokenHandler != nil {
						(*handler.PathTokenHandler).ServeHTTP(w, r, pathParam)
					} else if handler.ChainPathTokenHandler != nil {
						for _, value := range handler.ChainPathTokenHandler {
							if ok := value.ServeNextHTTP(w, r, pathParam); !ok {
								break
							}
						}
					} else {
						NotFound(w, r)
					}
				} else {
					NotFound(w, r)
				}
				return nil
			}
		}
	}
	return errors.New("cannot find match path param url " + r.URL.Path)
}

func handleUrl(c *config.Config, db *sql.DB, mux *http.ServeMux, w http.ResponseWriter, r *http.Request) {
	var err error
	//first try the path expression syntax
	err = handleUrlPathEx(c, db, mux, w, r)
	if err == nil {
		return
	}
	logUtil.DebugPrintln(err.Error())

	//second try the path param syntax
	err = handleUrlPathParam(c, db, mux, w, r)
	if err == nil {
		return
	}
	logUtil.DebugPrintln(err.Error())

	//third try the regular expression syntax
	err = handleUrlRegEx(c, db, mux, w, r)
	if err != nil {
		logUtil.DebugPrintln(err.Error())
		NotFound(w, r)
	}
}

func setupRootHandler(c *config.Config, db *sql.DB, mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", c.Site.Name)
		if r.URL.Path == "/" {
			io.WriteString(w, "I am alive!")
		} else {
			handleUrl(c, db, mux, w, r)
		}
	})
}

func setupStaticPath(c *config.Config, db *sql.DB, mux *http.ServeMux) {
	if stat, err := os.Stat(c.Site.StaticFilePath); err == nil && stat.IsDir() {
		fs := http.FileServer(http.Dir(c.Site.StaticFilePath))
		if fs != nil {
			lastEle := filepath.Base(c.Site.StaticFilePath)
			mux.Handle("/"+lastEle+"/", http.StripPrefix("/"+lastEle, fs))
		} else {
			log.Print("error getting handler for StaticFilePath")
		}
	} else {
		log.Printf("error access StaticFilePath: %v", err)
	}
}

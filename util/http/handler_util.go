package httpUtil

import (
	"database/sql"
	"log"
	"tiger/config"
)

// ENTRY POINT: perform any cleaning up of objects or anything else before server shutdown in here (if any)
func ShutdownCleanup(c *config.Config, db *sql.DB) {
	log.Print("shutdown cleanup ...")
	//////// add application specific logic below ////////
}

// ENTRY POINT: perform any pre-loading/caching of objects or anything else before server startup in here (if any)
func StartupInit(c *config.Config, db *sql.DB) {
	log.Print("startup init ...")
	//////// add application specific logic below ////////
}

// ENTRY POINT: register all custom error pages in here if not going to use default provided http.Error(w ResponseWriter, error string, code int) , http.NotFound(w ResponseWriter, r *Request)
// 	Example
// 	AddCustomErrorPage(http.StatusNotFound, "templates/errors/404Error.html", nil)
// 	AddCustomErrorPage(http.StatusNotFound, "templates/errors/404Error.html", map[string]string{ "custom header" : "can see?" })
// 	AddCustomErrorPage(http.StatusInternalServerError, "templates/errors/500Error.html", nil)
func RegisterCustomErrorPages(c *config.Config, db *sql.DB) {
	log.Print("register custom error pages ...")
	//////// add application specific logic below ////////
}

// ENTRY POINT: register all url mapping to handler in here
//
//call func AddHandler(urlMapping string, handler http.Handler, httpVerb ...string) where the valid values for httpVerb are as follow
//http.MethodGet     = "GET"
//http.MethodHead    = "HEAD"
//http.MethodPost    = "POST"
//http.MethodPut     = "PUT"
//http.MethodPatch   = "PATCH"
//http.MethodDelete  = "DELETE"
//http.MethodConnect = "CONNECT"
//http.MethodOptions = "OPTIONS"
//http.MethodTrace   = "TRACE"
//if not passed would default to http.MethodGet
//
//for support of Path Param url mapping like /{placeholder} or /:placeholder need to implement the httpUtil.PathTokenHandler interface before registering
//please call AddHandlerPathParam(urlMapping string, pathTokenHandler PathTokenHandler, httpVerb ...string)
//
//for support of more complicated url mapping like regular expression please call func AddHandlerRegEx(urlMapping string, handler http.Handler, httpVerb ...string) instead. Refer to go regexp package for the re syntax
//
//for support of chaining of handlers to call them one by one sequentially need to implement ChainNextHandler and/or ChainPathTokenHandler before registering
//please call their equivalent func AddChainHandler(...), AddChainHandlerRegEx(...), AddChainHandlerPathParam(...)
//
//for support of url rewriting please ensure the json attribute for UrlRewrite is set to true in config.json. due to performance concern this feature must be explicitly enabled. please call AddRewriteUrl(sourceUrl string, targetUrl string) where sourceUrl can be normal, path param, regular expression.
//for path param /{placeholder} or /:placeholder to be carried over to targetUrl ensure the SAME placeholder is placed in targetUrl.
//for regex matched to be carried over to targetUrl, please enclose in parenthesis on sourceUrl and then use $1 , $2 on targetUrl.
//
// 	Example
// 	AddHandlerRegEx("/hello1/.*/12[34]$", &logic1.ApiHandler{Db: db}, http.MethodGet)
// 	AddHandler("/hello2", &logic1.LogicHandler{Db: db}, http.MethodGet)
// 	AddHandler("/hello3", &logic2.ApiHandler{Config: c, Next: nil}, http.MethodGet)
// 	AddHandler("/hello4", &logic2.LogicHandler{}, http.MethodGet, http.MethodPost)
// 	AddHandlerPathParam("/hello5/:userId/test/{prodId}", &logic2.Logic2Handler{}, http.MethodGet, http.MethodPost)
//
// 	firstChain := []ChainNextHandler{
//		&logic3.Api1Handler{Config: c},
//		&logic3.Api2Handler{Config: c},
//	}
//	AddChainHandler("/hello6", firstChain, http.MethodGet)
//
//	secondChain := []ChainNextHandler{
//		&logic3.Api3Handler{Config: c},
//		&logic3.Api4Handler{Config: c},
//	}
//	AddChainHandlerRegEx("/hello7/.*/12[34]$", secondChain, http.MethodGet)
//
//	thirdChain := []ChainPathTokenHandler{
//		&logic3.Api5Handler{Config: c},
//		&logic3.Api6Handler{Config: c},
//	}
//	AddChainHandlerPathParam("/hello8/:userId/test/{prodId}", thirdChain, http.MethodGet)
//
//	fourthChain := []ChainNextHandler{
//		rateLimiter.NewTokenBucketHandler(60, 30, 30),
//		&logic3.Api7Handler{Config: c},
//	}
//	AddChainHandler("/hello9", fourthChain, http.MethodGet)
//
//	fifthChain := []ChainNextHandler{
//		rateLimiter.NewSlidingWindowHandler(60),
//		&logic3.Api7Handler{Config: c},
//	}
//	AddChainHandler("/hello10", fifthChain, http.MethodGet)
//
//	AddRewriteUrl("/testhello4", "/hello4")
//	AddRewriteUrl("/testhello5/haha/:userId/test/{prodId}", "/hello5/:userId/test/{prodId}")
//	AddRewriteUrl("/testhello1/haha/(.*)/(12[34]$)", "/hello1/$1/$2")
func RegisterHandler(c *config.Config, db *sql.DB) {
	log.Print("register handler ...")
	//////// add application specific logic below ////////
}
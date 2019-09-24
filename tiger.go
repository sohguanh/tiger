// main package where the http server will be started up/shutdown.
// 	How to use tiger framework
//
// 	Step 1
//	Ensure config/config.json are setup correctly for your environment.
//
// 	Step 2
// 	Start to add your application specific code in util/http/handler_util.go Refer to the relevant package documentation on how to do it.
//
// 	Step 3
// 	Compile by running go build. tiger.exe or tiger will be created.
//
// 	Step 4
// 	From Windows Command Prompt or Linux terminal, execute tiger.exe or tiger &
//
// 	Step 5
// 	Use a browser and navigate to your configured url in Step 1 config.json e.g http://localhost:8000
// 	You should see a message I am alive! This mean your http server is up and running.
// 	To shutdown, send a SIGINT signal. Ctrl-C for Windows Command Prompt. kill -SIGINT <pid> for Linux.
package main

import (	
	"flag"
	"log"	
	"strconv"	
	"net/http"
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"	
	"tiger/config"
	dbUtil "tiger/util/db"
	httpUtil "tiger/util/http"
	templateUtil "tiger/util/template"
	logUtil "tiger/util/log"
)

func main() {
	var env = ""
	env = os.Getenv("env")
	if env == "" {
		//try commandline option
		var flagVar string
		flag.StringVar(&flagVar, "env", "", "set environment setting to Dev,Qa,Prod")
		flag.Parse()
		env = flagVar
	}
	c, err := config.NewConfig(env)
	if err != nil { //cannot load config exit program
		panic(err)
	}
	
	if c.Site.LogToFile {
		f, err := config.NewLogFileName()
		if err == nil { //cannot set then print to stdout else print to file
			defer f.Close()
			log.SetOutput(f)
		}
	}	
	
	logUtil.SetLevel(logUtil.ValidLogLevel[strings.ToUpper(c.Site.LogLevel)])
	logUtil.DebugPrintf("%+v\n", c)
	
	db, err := dbUtil.NewDb(c)
	if err != nil { //cannot open db exit program
		panic(err)
	}
	defer db.Close()
			
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpUtil.StartupInit(c, db)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpUtil.RegisterCustomErrorPages(c, db)
	}()	
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpUtil.RegisterHandler(c, db)
	}()
	if c.TemplateConfig.Enable {
		wg.Add(1)
		go func() {
			defer wg.Done()
			templateUtil.NewTemplateUtil(c, db)
		}()
	}
	wg.Wait()
	
	actualMux := httpUtil.NewServeMux(c, db)
		
	var mux *http.ServeMux	
	if !c.Site.UrlRewrite {
		mux = actualMux
	} else {
		muxWrapper := http.NewServeMux()
		muxWrapper.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {			
								logUtil.DebugPrintln("muxWrapper incoming url: "+r.URL.Path)
								r.URL.Path = httpUtil.GetRewriteUrlTarget(r.URL.Path)
								logUtil.DebugPrintln("muxWrapper outgoing url: "+r.URL.Path)
								actualMux.ServeHTTP(w, r)
							})		
		mux = muxWrapper
	}

	connClosed := make(chan string)
	srv := &http.Server{
		Addr : ":"+strconv.Itoa(c.Site.Port), 
		Handler : mux,
	}
	srv.RegisterOnShutdown(func(){
		log.Print("received an interrupt signal, server shutting down ...")
		httpUtil.ShutdownCleanup(c, db)
		connClosed <- "server shutdown ..."
	})	
	go func(){
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint //block until interrupt signal is received
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Site.GracefulShutdownSec)*time.Second)
		defer cancel()		
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("error shutdown server: %v", err)
		}
	}()
	go func(){
		log.Print("server starting up ...")
		srv.ListenAndServe()
	}()
	go func(){
		client := &http.Client{Timeout: time.Duration(c.Site.CheckAliveTimeoutSec)*time.Second}	
		resp, err := client.Get("http://"+c.Site.Url+":"+strconv.Itoa(c.Site.Port))		
        if err != nil {
            log.Printf("error contact server: %v", err)
			return
        }
		defer resp.Body.Close()
        if resp.StatusCode == http.StatusOK {
            log.Print("server started up ...")
        }		
	}()
	
	log.Fatalln(<-connClosed) //block until message is received on the channel
}

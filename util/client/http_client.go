package clientUtil

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var httpClientPool = sync.Pool{
	New: func() interface{} {
		return &http.Client{}
	},
}

// HttpNewRequest. call Go http.NewRequest underlying
func HttpNewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

// HttpDo. call Go client.Do underlying but with the retry feature based on the passed in reqInfo parameter.
func HttpDo(r *http.Request, reqInfo *RequestInfo) ([]byte, error) {
	return httpDoCommon(r, reqInfo)
}

// HttpGet. call Go client.Get underlying but with the retry feature based on the passed in reqInfo parameter.
func HttpGet(url string, reqInfo *RequestInfo) ([]byte, error) {
	return httpCommon(url, reqInfo, http.MethodGet, "", nil, nil)
}

// HttpHead. call Go client.Head underlying but with the retry feature based on the passed in reqInfo parameter.
func HttpHead(url string, reqInfo *RequestInfo) ([]byte, error) {
	return httpCommon(url, reqInfo, http.MethodHead, "", nil, nil)
}

// HttpPost. call Go client.Post underlying but with the retry feature based on the passed in reqInfo parameter.
func HttpPost(url string, contentType string, body io.Reader, reqInfo *RequestInfo) ([]byte, error) {
	return httpCommon(url, reqInfo, http.MethodPost, contentType, body, nil)
}

// HttpPostForm. call Go client.PostForm underlying but with the retry feature based on the passed in reqInfo parameter.
func HttpPostForm(url string, data url.Values, reqInfo *RequestInfo) ([]byte, error) {
	return httpCommon(url, reqInfo, http.MethodPost+"Form", "", nil, data)
}

func httpCommon(url string, reqInfo *RequestInfo, httpVerb string, contentType string, postBody io.Reader, data url.Values) ([]byte, error) {
	if url == "" {
		return nil, errors.New("empty url")
	}
	client := httpClientPool.Get()
	defer client.(*http.Client).CloseIdleConnections()
	defer httpClientPool.Put(client)

	client.(*http.Client).Timeout = time.Duration(reqInfo.TimeoutSec) * time.Second

	var (
		r         *http.Response
		err       error
		retryCnt  = 1
		bodyBytes []byte
	)
	r, err = httpCommonRetry(client.(*http.Client), url, reqInfo.TimeoutSec, httpVerb, contentType, postBody, data)
	for err != nil && retryCnt < reqInfo.RetryTimes {
		retryCnt++
		<-time.After(time.Duration(reqInfo.WaitBeforeRetrySec) * time.Second)
		r, err = httpCommonRetry(client.(*http.Client), url, reqInfo.TimeoutSec, httpVerb, contentType, postBody, data)
	}
	if err == nil {
		bodyBytes, err = ioutil.ReadAll(r.Body)
		r.Body.Close()
		return bodyBytes, err
	} else {
		return nil, err
	}
}

func httpCommonRetry(client *http.Client, url string, timeoutSec int, httpVerb string, contentType string, postBody io.Reader, data url.Values) (*http.Response, error) {
	var (
		r   *http.Response
		err error
	)
	if httpVerb == http.MethodGet {
		r, err = client.Get(url)
	} else if httpVerb == http.MethodPost {
		r, err = client.Post(url, contentType, postBody)
	} else if httpVerb == http.MethodPost+"Form" {
		r, err = client.PostForm(url, data)
	} else if httpVerb == http.MethodHead {
		r, err = client.Head(url)
	}
	return r, err
}

func httpDoCommon(req *http.Request, reqInfo *RequestInfo) ([]byte, error) {
	client := httpClientPool.Get()
	defer client.(*http.Client).CloseIdleConnections()
	defer httpClientPool.Put(client)

	client.(*http.Client).Timeout = time.Duration(reqInfo.TimeoutSec) * time.Second

	var (
		r         *http.Response
		err       error
		retryCnt  = 1
		bodyBytes []byte
	)
	r, err = client.(*http.Client).Do(req)
	for err != nil && retryCnt < reqInfo.RetryTimes {
		retryCnt++
		<-time.After(time.Duration(reqInfo.WaitBeforeRetrySec) * time.Second)
		r, err = client.(*http.Client).Do(req)
	}
	if err == nil {
		bodyBytes, err = ioutil.ReadAll(r.Body)
		r.Body.Close()
		return bodyBytes, err
	} else {
		return nil, err
	}
}

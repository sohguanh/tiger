package httpUtil

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var onceRewriteUrl sync.Once
var mapRewriteUrl map[string]string
var mutexRewriteUrl sync.RWMutex

// AddRewriteUrl. sourceUrl parameter placeholder syntax is () targetUrl parameter substituition syntax is $1 $ 2 etc.
// 	Example sourceUrl /(id)  and targetUrl /$1
func AddRewriteUrl(sourceUrl string, targetUrl string) {
	initRewriteUrl()
	mutexRewriteUrl.Lock()
	defer mutexRewriteUrl.Unlock()
	mapRewriteUrl[strings.TrimSpace(sourceUrl)] = strings.TrimSpace(targetUrl)
}

// GetRewriteUrlTarget to get the target rewritten url based on the sourceUrl parameter.
func GetRewriteUrlTarget(sourceUrl string) string {
	initRewriteUrl()
	mutexRewriteUrl.RLock()
	defer mutexRewriteUrl.RUnlock()
	sourceUrl = strings.TrimSpace(sourceUrl)
	for src, tgt := range mapRewriteUrl {
		if ok, url := matchRewriteUrlSource(sourceUrl, src, tgt); ok {
			return url
		}
	}
	return sourceUrl
}

func matchRewriteUrlSource(incomingSourceUrl, mapSourceUrl, mapTargetUrl string) (bool, string) {
	var (
		found bool
	)
	if mapSourceUrl == incomingSourceUrl { //direct match
		return true, mapTargetUrl
	}
	if found, _ = filepath.Match(mapSourceUrl, incomingSourceUrl); found { //filepath match
		return true, mapTargetUrl
	}

	actualToken := splitBySlashToken(incomingSourceUrl)
	mapSourceToken := splitBySlashToken(mapSourceUrl)
	var pathParam = make(map[string]string)
	found = false
	if len(mapSourceToken) == len(actualToken) { //path param match
		for index, value := range mapSourceToken {
			if matched := pathParamRE.MatchString(value); matched {
				pathParam[value] = actualToken[index]
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
		for placeholder, value := range pathParam {
			mapTargetUrl = strings.ReplaceAll(mapTargetUrl, placeholder, value)
		}
		return true, mapTargetUrl
	}

	if srcRe, err := regexp.Compile(mapSourceUrl); err == nil {
		if matched := srcRe.FindStringSubmatch(incomingSourceUrl); matched != nil {
			for index, value := range matched {
				mapTargetUrl = strings.ReplaceAll(mapTargetUrl, "$"+strconv.Itoa(index), value)
			}
			return true, mapTargetUrl
		}
	}

	return false, ""
}

func initRewriteUrl() {
	onceRewriteUrl.Do(func() { //singleton
		mapRewriteUrl = make(map[string]string)
	})
}

// i18nUtil package has dependency on a third party package. please install it first before using this package.
//
// 	go get -v golang.org/x/text
package i18nUtil

import (
	"bufio"
	"errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"os"
	"strings"
	"sync"
	logUtil "tiger/util/log"
)

type tagInfo struct {
	LangTag     language.Tag
	LangPrinter *message.Printer
}

var onceTag sync.Once
var mutexTag sync.RWMutex

var mapTag map[string]tagInfo

const (
	keySep   string = "|"
	equalSep string = "="
)

// GetMsg to get the actual language text.
// tag parameter come from https://godoc.org/golang.org/x/text/language/display
// 	English en
// 	AmericanEnglish en-US
// 	BritishEnglish en-GB
// 	Chinese zh
// 	SimplifiedChinese zh-Hans
// 	TraditionalChinese zh-Hant
//if key parameter value inside the properties file has substituition placeholder like %s or %d or %f please pass the actual run-time value into the param ...interface{}
func GetMsg(tag string, key string, param ...interface{}) string {
	initTagHandler()
	mutexTag.RLock()
	defer mutexTag.RUnlock()
	if value, found := mapTag[tag]; found {
		if param == nil {
			return value.LangPrinter.Sprintf(message.Key(key+keySep+tag, ""), nil)
		} else {
			return value.LangPrinter.Sprintf(message.Key(key+keySep+tag, ""), param...)
		}
	}
	return ""
}

// LoadProperties to parse the properties file to be stored internally.
// tag parameter values please refer to GetMsg
func LoadProperties(tag string, propsFilename string) error {
	initTagHandler()
	mutexTag.Lock()
	defer mutexTag.Unlock()

	if _, err := os.Stat(propsFilename); os.IsNotExist(err) {
		return errors.New("cannot find properties file " + propsFilename)
	}

	file, err := os.Open(propsFilename)
	if err != nil {
		return errors.New("cannot open properties file " + propsFilename)
	}
	defer file.Close()
	logUtil.DebugPrint("process " + propsFilename)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, equalSep); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				addToTag(tag, key, value)
			}
		}
	}
	return nil
}

func addToTag(tag string, key string, msg string) {
	if _, found := mapTag[tag]; !found {
		langTag := message.MatchLanguage(tag)
		langPrinter := message.NewPrinter(langTag)
		mapTag[tag] = tagInfo{LangTag: langTag, LangPrinter: langPrinter}
	}
	//not sure is library bug or not but same key value for different tag is overwriting each other
	//to ensure unique-ness append with keySep and tag
	err := message.SetString(mapTag[tag].LangTag, key+keySep+tag, msg)
	if err != nil {
		log.Print("cannot set key " + key + " msg " + msg)
	}
}

func initTagHandler() {
	onceTag.Do(func() { //singleton
		mapTag = make(map[string]tagInfo)
	})
}

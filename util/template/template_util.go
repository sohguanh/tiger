// templateUtil is the package that cache the template.Template object so it can be reused across different calls from callers.
package templateUtil

import (
	"database/sql"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"tiger/config"
	logUtil "tiger/util/log"
)

var onceTemplate sync.Once
var onceHandler sync.Once
var mutexMapTemplate sync.RWMutex
var mapTemplate map[string]*template.Template

// GetTemplate to retrieve the template.Template object.
// template parameter is the full path to the template file where path separator are set to /
func GetTemplate(template string) (*template.Template, error) {
	initMapHandler()
	mutexMapTemplate.RLock()
	defer mutexMapTemplate.RUnlock()
	if value, found := mapTemplate[template]; found {
		return value, nil
	} else {
		return nil, errors.New("cannot find " + template)
	}
}

// AddTemplate to add and parse the template into template.Template object to be reused when call by GetTemplate.
// template parameter is the full path to the template file where path separator are set to /
func AddTemplate(template string) (*template.Template, error) {
	initMapHandler()
	mutexMapTemplate.Lock()
	defer mutexMapTemplate.Unlock()
	tpl, err := addToMapHandler(template)
	return tpl, err
}

// ListTemplate to list down all the template parameters stored internally. To check if the template has been added previously.
func ListTemplate() {
	initMapHandler()
	mutexMapTemplate.RLock()
	defer mutexMapTemplate.RUnlock()
	for key, _ := range mapTemplate {
		log.Println(key)
	}
}

// NewTemplateUtil is to walk recursively through the template folder and parse all the template files into template.Template objects. The configuration are from the json attribute called TemplateConfig in config.json.
func NewTemplateUtil(c *config.Config, db *sql.DB) {
	onceTemplate.Do(func() { //singleton
		logUtil.DebugPrint("template first time init\n")
		initMapHandler()
		root := filepath.Base(c.TemplateConfig.Path)

		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if info != nil && !info.IsDir() && strings.HasSuffix(path, c.TemplateConfig.FileExt) {
				_, err := addToMapHandler(path)
				if err != nil {
					return err
				}
			}
			return nil
		})
	})
}

func addToMapHandler(path string) (*template.Template, error) {
	slashPath := filepath.ToSlash(path)
	logUtil.DebugPrint("process " + slashPath)

	b, err := ioutil.ReadFile(slashPath)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New(slashPath).Parse(string(b))
	if err != nil {
		return nil, err
	}
	mapTemplate[slashPath] = tpl
	return tpl, nil
}

func initMapHandler() {
	onceHandler.Do(func() { //singleton
		mapTemplate = make(map[string]*template.Template)
	})
}

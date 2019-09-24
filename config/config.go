// config is the package that is doing the parsing of config.json into a Config object to be used for the application. It also include the creation of a new logfile.
package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// EnvProd , EnvDev point to the json attribute that is defined in config.json
// It is highly recommended not to change this but if you know what you are doing, feel free to change it.
const (
	EnvProd = `Prod`
	EnvDev  = `Dev`
)

// ConfigFileName is the filename where all configuration are stored.
// Change this to another filename if you want.
var ConfigFileName = `config\config.json`

// LogFileName is the filename of the log.
// Change this to another filename if you want.
var LogFileName = `tiger.log`

// Config is the struct that contain all the configuration that is from config.json
type Config struct {
	Env  string
	Site struct {
		Name                 string
		Url                  string
		Port                 int
		LogToFile            bool
		LogLevel             string
		GracefulShutdownSec  int
		CheckAliveTimeoutSec int
		ReadTimeoutSec       int
		ReadHeaderTimeoutSec int
		WriteTimeoutSec      int
		IdleTimeoutSec       int
		MaxHeaderBytes       int
		StaticFilePath       string
		UrlRewrite           bool
	}
	Database struct {
		Name     string
		Host     string
		Port     int
		Username string
		Password string
	}
	TemplateConfig struct {
		Enable  bool
		Path    string
		FileExt string
	}
}

type fileConfig struct {
	Env string `json:"Env"`
	Dev struct {
		Site struct {
			Name                 string `json:"Name"`
			Url                  string `json:"Url"`
			Port                 int    `json:"Port"`
			LogToFile            bool   `json:"LogToFile"`
			LogLevel             string `json:"LogLevel"`
			GracefulShutdownSec  int    `json:"GracefulShutdownSec"`
			CheckAliveTimeoutSec int    `json:"CheckAliveTimeoutSec"`
			ReadTimeoutSec       int    `json:"ReadTimeoutSec"`
			ReadHeaderTimeoutSec int    `json:"ReadHeaderTimeoutSec"`
			WriteTimeoutSec      int    `json:"WriteTimeoutSec"`
			IdleTimeoutSec       int    `json:"IdleTimeoutSec"`
			MaxHeaderBytes       int    `json:"MaxHeaderBytes"`
			StaticFilePath       string `json:"StaticFilePath"`
			UrlRewrite           bool   `json:"UrlRewrite"`
		}
		Database struct {
			Name     string `json:"Name"`
			Host     string `json:"Host"`
			Port     int    `json:"Port"`
			Username string `json:"Username"`
			Password string `json:"Password"`
		}
		TemplateConfig struct {
			Enable  bool   `json:"Enable"`
			Path    string `json:"Path"`
			FileExt string `json:"FileExt"`
		}
	}
	Prod struct {
		Site struct {
			Name                 string `json:"Name"`
			Url                  string `json:"Url"`
			Port                 int    `json:"Port"`
			LogToFile            bool   `json:"LogToFile"`
			LogLevel             string `json:"LogLevel"`
			GracefulShutdownSec  int    `json:"GracefulShutdownSec"`
			CheckAliveTimeoutSec int    `json:"CheckAliveTimeoutSec"`
			ReadTimeoutSec       int    `json:"ReadTimeoutSec"`
			ReadHeaderTimeoutSec int    `json:"ReadHeaderTimeoutSec"`
			WriteTimeoutSec      int    `json:"WriteTimeoutSec"`
			IdleTimeoutSec       int    `json:"IdleTimeoutSec"`
			MaxHeaderBytes       int    `json:"MaxHeaderBytes"`
			StaticFilePath       string `json:"StaticFilePath"`
			UrlRewrite           bool   `json:"UrlRewrite"`
		}
		Database struct {
			Name     string `json:"Name"`
			Host     string `json:"Host"`
			Port     int    `json:"Port"`
			Username string `json:"Username"`
			Password string `json:"Password"`
		}
		TemplateConfig struct {
			Enable  bool   `json:"Enable"`
			Path    string `json:"Path"`
			FileExt string `json:"FileExt"`
		}
	}
}

var config fileConfig
var retnConfig Config
var once sync.Once
var configErr bool = false

// NewConfig is to get a singleton Config object from the package.
//
// env parameter is passed in from caller. caller will attempt to get env based on below steps.
// 	Step 1 attempt to read the value from environment variable called env
// 	Step 2 if Step 1 fail, attempt to read from command-line option called env
// 	Step 3 if step 2 fail, the env parameter will be ""
//	Step 4 if step 3 return as "", default to read from the json attribute called Env in config.json
//
//	the valid values for env is Dev , Prod
//	any other environment please amend config.json and config.go accordingly
//	to make adding new environment more flexible, please send suggestion to sohguanh@gmail.com
func NewConfig(env string) (*Config, error) {
	once.Do(func() { //singleton
		dir, err := os.Getwd()
		f, err := ioutil.ReadFile(dir + string(os.PathSeparator) + ConfigFileName)
		if err != nil {
			log.Fatal(err)
			configErr = true
		}
		config = fileConfig{}
		err = json.Unmarshal(f, &config)
		if err != nil {
			log.Fatal(err)
			configErr = true
		}
		var configEnv string
		if env == "" {
			configEnv = config.Env
		} else {
			configEnv = env
		}
		if configEnv == EnvProd {
			retnConfig.Env = config.Env

			retnConfig.Site.Name = config.Prod.Site.Name
			retnConfig.Site.Url = config.Prod.Site.Url
			retnConfig.Site.Port = config.Prod.Site.Port
			retnConfig.Site.LogToFile = config.Prod.Site.LogToFile
			retnConfig.Site.LogLevel = config.Prod.Site.LogLevel
			retnConfig.Site.GracefulShutdownSec = config.Prod.Site.GracefulShutdownSec
			retnConfig.Site.CheckAliveTimeoutSec = config.Prod.Site.CheckAliveTimeoutSec
			retnConfig.Site.ReadTimeoutSec = config.Prod.Site.ReadTimeoutSec
			retnConfig.Site.ReadHeaderTimeoutSec = config.Prod.Site.ReadHeaderTimeoutSec
			retnConfig.Site.WriteTimeoutSec = config.Prod.Site.WriteTimeoutSec
			retnConfig.Site.IdleTimeoutSec = config.Prod.Site.IdleTimeoutSec
			retnConfig.Site.MaxHeaderBytes = config.Prod.Site.MaxHeaderBytes
			retnConfig.Site.StaticFilePath = config.Prod.Site.StaticFilePath
			retnConfig.Site.UrlRewrite = config.Prod.Site.UrlRewrite

			retnConfig.Database.Name = config.Prod.Database.Name
			retnConfig.Database.Host = config.Prod.Database.Host
			retnConfig.Database.Port = config.Prod.Database.Port
			retnConfig.Database.Username = config.Prod.Database.Username
			retnConfig.Database.Password = config.Prod.Database.Password

			retnConfig.TemplateConfig.Enable = config.Prod.TemplateConfig.Enable
			retnConfig.TemplateConfig.Path = config.Prod.TemplateConfig.Path
			retnConfig.TemplateConfig.FileExt = config.Prod.TemplateConfig.FileExt
		} else { //default to Dev
			retnConfig.Env = config.Env

			retnConfig.Site.Name = config.Dev.Site.Name
			retnConfig.Site.Url = config.Dev.Site.Url
			retnConfig.Site.Port = config.Dev.Site.Port
			retnConfig.Site.LogToFile = config.Dev.Site.LogToFile
			retnConfig.Site.LogLevel = config.Dev.Site.LogLevel
			retnConfig.Site.GracefulShutdownSec = config.Dev.Site.GracefulShutdownSec
			retnConfig.Site.CheckAliveTimeoutSec = config.Dev.Site.CheckAliveTimeoutSec
			retnConfig.Site.ReadTimeoutSec = config.Dev.Site.ReadTimeoutSec
			retnConfig.Site.ReadHeaderTimeoutSec = config.Dev.Site.ReadHeaderTimeoutSec
			retnConfig.Site.WriteTimeoutSec = config.Dev.Site.WriteTimeoutSec
			retnConfig.Site.IdleTimeoutSec = config.Dev.Site.IdleTimeoutSec
			retnConfig.Site.MaxHeaderBytes = config.Dev.Site.MaxHeaderBytes
			retnConfig.Site.StaticFilePath = config.Dev.Site.StaticFilePath
			retnConfig.Site.UrlRewrite = config.Dev.Site.UrlRewrite

			retnConfig.Database.Name = config.Dev.Database.Name
			retnConfig.Database.Host = config.Dev.Database.Host
			retnConfig.Database.Port = config.Dev.Database.Port
			retnConfig.Database.Username = config.Dev.Database.Username
			retnConfig.Database.Password = config.Dev.Database.Password

			retnConfig.TemplateConfig.Enable = config.Dev.TemplateConfig.Enable
			retnConfig.TemplateConfig.Path = config.Dev.TemplateConfig.Path
			retnConfig.TemplateConfig.FileExt = config.Dev.TemplateConfig.FileExt
		}
	})
	if !configErr {
		return &retnConfig, nil
	} else {
		return nil, errors.New("config error")
	}
}

// NewLogFileName is to get a file that is created based on the LogFileName variable.
func NewLogFileName() (*os.File, error) {
	dir, err := os.Getwd()
	f, err := os.OpenFile(dir+string(os.PathSeparator)+LogFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil { //no need exit program if cannot open file as default log write to stdout
		log.Print(err)
		return nil, err
	}
	//defer f.Close()
	return f, nil
}

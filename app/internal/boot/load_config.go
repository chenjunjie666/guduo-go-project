package boot

import (
	"fmt"
	"github.com/gobuffalo/packr"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type AllConfig struct {
	Database *struct {
		Crawler *struct {
			Host   string
			User   string
			Pass   string
			Dbname string
			Port   int
		}
		Clean *struct {
			Host   string
			User   string
			Pass   string
			Dbname string
			Port   int
		}
		Lolipop *struct {
			Host   string
			User   string
			Pass   string
			Dbname string
			Port   int
		}
		Carl *struct {
			Host   string
			User   string
			Pass   string
			Dbname string
			Port   int
		}
	}
	Redis *struct {
		Host string
		Pass string
		Port int
		DB   int
	}
	Proxy *struct {
		Host   string
		Port   string
		Secret string
	}
}

var Cfg *AllConfig

func init() {
	yamlName := "config.yml"

	cfgPath := "../boot"
	box := packr.NewBox(cfgPath)
	yamlFile, e := box.Find(yamlName)

	if e != nil {
		yamlFile, e = ioutil.ReadFile(yamlName)
		if e != nil {
			ex, err := os.Executable()
			if err != nil {
				panic(err)
			}
			yamlFile, e = ioutil.ReadFile(filepath.Dir(ex) + "/config.yml")
			if e != nil {
				log.Panic(fmt.Sprintf("读取配置文件出错，错误详情：%s", e))
				panic(fmt.Sprintf("读取配置文件出错，请检查配置，错误详情：%s", e))
			}
		}
	}

	e = yaml.Unmarshal(yamlFile, &Cfg)
	if e != nil {
		log.Panic(fmt.Sprintf("解析配置文件出错，错误详情：%s", e))
		panic(fmt.Sprintf("解析配置文件出错，请检查配置，错误详情：%s", e))
	}
}

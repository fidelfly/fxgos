package system

import (
	"os"
	"github.com/sirupsen/logrus"
	"github.com/BurntSushi/toml"
	"flag"
)

const DEF_CONFIG_FILE = "config.toml"

func InitConfig() (err error) {
	// Parse Command Parameters
	var configFile = ""
	flag.StringVar(&configFile, "config", DEF_CONFIG_FILE, "Set Config File")
	flag.Parse()

	// Parse Config File
	err = ParseConfig(configFile, &Properties)
	return
}

func ParseConfig(file string, target interface{}) (err error) {
	if _, err = os.Stat(file); err == nil {
		logrus.Infof("Config file : %s Found!", file)
		_, err = toml.DecodeFile(file, target)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		logrus.Errorf("Config file : %s is not found!", file)
	}
	return
}

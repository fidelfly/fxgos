package fxgos

import (
	"os"
	"github.com/sirupsen/logrus"
	"github.com/BurntSushi/toml"
)

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

package app

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"

	"github.com/fidelfly/fxgos/system"
)

func initDatabase(config system.DatabaseProperties) (err error) {
	system.DbEngine, err = xorm.NewEngine("mysql", getDBUrl(config))
	if err == nil {
		err = system.DbEngine.Ping()
		if err == nil {
			logrus.Info("Database is connected!")
		} else {
			panic(err)
		}
	}
	return
}

func getDBUrl(config system.DatabaseProperties) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", config.User, config.Password, config.Host, config.Port, config.Schema)
}

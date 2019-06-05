package main

import (
	"fmt"

	"github.com/fidelfly/fxgo/logx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	"github.com/fidelfly/fxgos/system"
)

func initDatabase(config *system.DatabaseProperties) (err error) {
	system.DbEngine, err = xorm.NewEngine("mysql", getDBUrl(config))
	if err == nil {
		err = system.DbEngine.Ping()
		if err == nil {
			logx.Info("Database is connected!")
		} else {
			logx.Errorf("Can't connect to database")
		}
	}
	return
}

func getDBUrl(config *system.DatabaseProperties) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", config.User, config.Password, config.Host, config.Port, config.Schema)
}

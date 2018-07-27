package system

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	"fmt"
)

var DbEngine *xorm.Engine

func InitDatabase(config DatabaseProperties) (err error){
	DbEngine, err = xorm.NewEngine("mysql", getDBUrl(config))
	if err == nil {
		err = DbEngine.Ping()
		if err != nil {
			logrus.Info("Database is connected!")
		}
	}
	return
}

func getDBUrl(config DatabaseProperties) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", config.User, config.Password, config.Host, config.Port, config.Schema)
}


package db

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	xormc "xorm.io/core"

	"github.com/fidelfly/gox/logx"
)

type Server struct {
	Host     string
	Port     int64
	Schema   string
	User     string
	Password string
	Params   map[string]interface{}
}

func (db Server) getUrl() string {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db.User, db.Password, db.Host, db.Port, db.Schema)
	if len(db.Params) > 0 {
		params := make([]string, 0)
		for key, val := range db.Params {
			if len(key) > 0 && val != nil {
				params = append(params, fmt.Sprintf("%s=%v", key, val))
			}
		}
		if len(params) > 0 {
			url += "?" + strings.Join(params, "&")
		}
	}
	//return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=Local", db.User, db.Password, db.Host, db.Port, db.Schema)
	return url
}

func (db Server) getTarget() string {
	return fmt.Sprintf("%s:%d", db.Host, db.Port)
}

type Logger struct {
	*logx.Logger
	level   xormc.LogLevel
	showSql bool
}

func (dl *Logger) ShowSQL(show ...bool) {
	if len(show) > 0 {
		dl.showSql = show[0]
	} else {
		dl.showSql = true
	}
}

func (dl *Logger) IsShowSQL() bool {
	return dl.showSql
}

func (dl *Logger) Level() xormc.LogLevel {
	return dl.level
}
func (dl *Logger) SetLevel(l xormc.LogLevel) {
	dl.level = l
	switch dl.level {
	case xormc.LOG_DEBUG:
		dl.Logger.Level = logrus.DebugLevel
	case xormc.LOG_INFO:
		dl.Logger.Level = logrus.InfoLevel
	case xormc.LOG_WARNING:
		dl.Logger.Level = logrus.WarnLevel
	case xormc.LOG_ERR:
		dl.Logger.Level = logrus.ErrorLevel
	case xormc.LOG_OFF:
		dl.Logger.Level = logrus.PanicLevel
	}
}

func GetLogLevel(level string) xormc.LogLevel {
	switch level {
	case "info":
		return xormc.LOG_INFO
	case "debug":
		return xormc.LOG_DEBUG
	case "warning":
		return xormc.LOG_WARNING
	case "error":
		return xormc.LOG_ERR
	default:
		return xormc.LOG_WARNING
	}
}

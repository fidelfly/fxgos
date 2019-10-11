package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	"github.com/fidelfly/fxgo/logx"
)

//export
func InitEngine(instance *Server, options ...EngineOption) (err error) {
	engine, err := NewEngine(instance)
	if err == nil {
		Engine = engine
		SetEngineOption(options...)
	}
	return
}

func SetEngineOption(options ...EngineOption) {
	if Engine != nil {
		for _, opt := range options {
			opt(Engine)
		}
	}
}

type EngineOption func(engine *xorm.Engine)

//export
func NewEngine(instance *Server) (engine *xorm.Engine, err error) {
	engine, err = xorm.NewEngine("mysql", instance.getUrl())
	if err == nil {
		terr := engine.Ping()
		if terr == nil {
			logx.Infof("Database(%s) is available!", instance.getTarget())
		} else {
			logx.Errorf("Can't connect to database(%s)", instance.getTarget())
		}
	}
	return
}

package model

import (
	"sync"

	"github.com/casbin/casbin"
	"github.com/casbin/casbin/model"
	"github.com/fidelfly/gox/logx"

	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/iam/res"
)

var enforcerCache sync.Map // map[string]*casbin.Enforcer

//export
func GetIAMEnforcer(resType string) *casbin.Enforcer {
	if e, ok := enforcerCache.Load(resType); ok {
		return e.(*casbin.Enforcer)
	}
	e := newIAMEnforcer(resType)
	if e == nil {
		return nil
	}
	ce, _ := enforcerCache.LoadOrStore(resType, e)
	return ce.(*casbin.Enforcer)
}

func newIAMEnforcer(resType string) *casbin.Enforcer {
	m := &res.Model{
		ResourceType: resType,
	}
	var iamModel model.Model
	if ok, err := db.Read(m); ok {
		iamModel = casbin.NewModel(string(m.Data))
	} else if err != nil {
		logx.Error(err)
		return nil
	}
	e := casbin.NewEnforcer(iamModel, NewAdapter(resType, m.Policy))

	return e
}

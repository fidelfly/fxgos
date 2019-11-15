package da

import (
	"github.com/casbin/casbin"
)

var enforcer *casbin.Enforcer

const dataModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, "sa") || sgMatch(r.sub, r.obj)
`

func initDaModel() {
	enforcer = casbin.NewEnforcer(casbin.NewModel(dataModel), newAdapter())
	enforcer.AddFunction("sgMatch", _sgMatchFunc)
}

type resource struct {
	Type string
	Id   int64
}

func _sgMatch(sub string, res resource) bool {
	usg := &userSg{}
	if err := sgCache.GetObject(sub, usg); err != nil {
		return false
	}
	rsg := &resourceSg{}
	if err := sgCache.GetObject(encodeResKey(res.Type, res.Id), rsg); err != nil {
		return false
	}
	for _, ug := range usg.SecurityGroups {
		for _, rg := range rsg.SecurityGroups {
			if ug == rg {
				return true
			}
		}
	}
	return false
}

func _sgMatchFunc(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return false, nil
	}
	if sub, ok := args[0].(string); ok {
		if res, ok := args[1].(resource); ok {
			return _sgMatch(sub, res), nil
		}
	}
	return false, nil
}

func updateDaModelByUser(userId int64, superAdmin bool) {
	sub := encodeUserKey(userId)
	enforcer.DeleteRolesForUser(sub)
	if superAdmin {
		enforcer.AddRoleForUser(sub, "sa")
	}
}

func removeUserFromDaModel(userId int64) {
	sub := encodeUserKey(userId)
	enforcer.DeleteRolesForUser(sub)
}

func validateDataAccess(sub string, res resource) bool {
	return enforcer.Enforce(sub, res)
}

package model

import (
	"github.com/casbin/casbin"

	"github.com/fidelfly/fxgos/cmd/service/api/iam"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/gostool/db"
)

func Validate(resType, sub, obj, act string) bool {
	e := GetIAMEnforcer(resType)
	if e == nil {
		return true
	}
	e.DeleteUser("abc")
	return e.Enforce(sub, obj, act)
}

//export
func UpdateRolePolicy(roleId int64, inheritRoles ...int64) {
	enforcerCache.Range(func(key, value interface{}) bool {
		resType := key.(string)
		e := value.(*casbin.Enforcer)
		roleSub := iam.EncodeRoleSubject(roleId)
		e.DeleteRolesForUser(roleSub)
		if len(inheritRoles) > 0 {
			for _, iRole := range inheritRoles {
				e.AddRoleForUser(roleSub, iam.EncodeRoleSubject(iRole))
			}
		}
		e.DeletePermissionsForUser(roleSub)
		policies := make([]*res.Policy, 0)
		err := db.Find(&policies, db.Where("resource_type = ? and role_id = ?", resType, roleId))
		if err == nil && len(policies) > 0 {
			for _, p := range policies {
				if len(p.Act) > 0 {
					for _, act := range p.Act {
						e.AddPolicy(p.Sub, p.Obj, act)
					}
				}
			}
		}
		return true
	})
}

func DeleteRolePolicy(roleId int64) {
	enforcerCache.Range(func(key, value interface{}) bool {
		e := value.(*casbin.Enforcer)
		roleSub := iam.EncodeRoleSubject(roleId)
		e.DeleteRole(roleSub)
		e.DeletePermissionsForUser(roleSub)
		return true
	})
}

func DeleteUserPolicy(userId int64) {
	enforcerCache.Range(func(key, value interface{}) bool {
		e := value.(*casbin.Enforcer)
		userSub := iam.EncodeUserSubject(userId)
		e.DeletePermissionsForUser(userSub)
		e.DeleteUser(userSub)
		return true
	})
}

//export
func UpdateUserRole(userId int64, roles []int64, superAdmin bool) {
	enforcerCache.Range(func(key, value interface{}) bool {
		//resType := key.(string)
		e := value.(*casbin.Enforcer)
		e.DeleteRolesForUser(iam.EncodeUserSubject(userId))
		//e.DeleteUser(iamx.EncodeUserSubject(userId))
		if superAdmin {
			e.AddRoleForUser(iam.EncodeUserSubject(userId), "sa")
		}
		for _, role := range roles {
			e.AddRoleForUser(iam.EncodeUserSubject(userId), iam.EncodeRoleSubject(role))
		}
		return true
	})
}

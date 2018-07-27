package caches

import (
	"strconv"
	"github.com/lyismydg/fxgos/system"
	"github.com/patrickmn/go-cache"
)

type UserInfo struct {
	Id int64
	Code string
	TenantId int64
	TenantName string
	TenantCode string
}

func userInfoResolver(key string) interface{} {
	userId, _ := strconv.ParseInt(key, 10, 64)
	userInfo := new(UserInfo)
	system.DbEngine.SQL("select a.id, a.code, a.tenant_id, b.code as tenant_code, b.name as tenant_name from user as a, tenant as b where a.id = ? and a.tenant_id = b.id", userId).Get(userInfo)
	return userInfo
}

func init() {
	system.UserCache = system.CreateEnsureCache(cache.NoExpiration, 0, userInfoResolver)
}


package caches

import (
	"strconv"
	"github.com/lyismydg/fxgos/system"
	"github.com/patrickmn/go-cache"
)

type UserInfo struct {
	Id int64
	Code string
}

func userInfoResolver(key string) interface{} {
	userId, _ := strconv.ParseInt(key, 10, 64)
	userInfo := new(UserInfo)
	system.DbEngine.SQL("select a.id, a.code from user as a where a.id = ?", userId).Get(userInfo)
	return userInfo
}

func init() {
	system.UserCache = system.CreateEnsureCache(cache.NoExpiration, 0, userInfoResolver)
}


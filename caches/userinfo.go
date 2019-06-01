package caches

import (
	"strconv"

	"github.com/fidelfly/fxgo/cachex"
	"github.com/fidelfly/fxgo/logx"

	"github.com/fidelfly/fxgos/system"
)

type UserInfo struct {
	Id   int64
	Code string
}

func userInfoResolver(key string) interface{} {
	userId, _ := strconv.ParseInt(key, 10, 64)
	userInfo := new(UserInfo)
	_, err := system.DbEngine.SQL("select a.id, a.code from user as a where a.id = ?", userId).Get(userInfo)
	if err != nil {
		logx.Errorf("error found during reading user(=%s) information from database for cache", key)
	}
	return userInfo
}

func init() {
	system.UserCache = cachex.CreateEnsureCache(cachex.NoExpiration, 0, userInfoResolver)
}

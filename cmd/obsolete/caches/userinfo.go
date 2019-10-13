package caches

import (
	"strconv"

	"github.com/fidelfly/gox/cachex/mcache"
	"github.com/fidelfly/gox/logx"

	"github.com/fidelfly/fxgos/cmd/utilities/system"
)

type UserInfo struct {
	ID   int64
	Code string
}

func userInfoResolver(key string) interface{} {
	userID, _ := strconv.ParseInt(key, 10, 64)
	userInfo := new(UserInfo)
	_, err := system.DbEngine.SQL("select a.id, a.code from user as a where a.id = ?", userID).Get(userInfo)
	if err != nil {
		logx.Errorf("error found during reading user(=%s) information from database for cache", key)
	}
	return userInfo
}

func init() {
	system.UserCache = mcache.NewEnsureCache(mcache.NoExpiration, 0, userInfoResolver)
}

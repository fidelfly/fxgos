package user

import (
	"context"
	"strconv"
	"time"

	"github.com/fidelfly/gox/cachex/mcache"

	"github.com/fidelfly/fxgos/cmd/utilities/pub"
)

type CacheInfo struct {
	Id    int64
	Code  string
	Name  string
	Email string
}

var myCache = mcache.NewEnsureCache(mcache.DefaultExpiration, 5*time.Minute, userResolver)

func userResolver(key string) interface{} {
	id, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return nil
	}
	u, err := Read(context.Background(), id)
	if err != nil {
		return nil
	}

	return &CacheInfo{
		Id:    u.Id,
		Code:  u.Code,
		Email: u.Email,
		Name:  u.Name,
	}
}

func cacheSubscriber(pubData interface{}) error {
	if re, ok := pubData.(pub.ResourceEvent); ok {
		if re.Type == ResourceType && re.Action != pub.ResourceCreate {
			myCache.Remove(strconv.FormatInt(re.Id, 10))
		}
	}
	return nil
}

func GetCache(userId int64) *CacheInfo {
	if userId > 0 {
		if cacheObj, ok := myCache.Get(strconv.FormatInt(userId, 10)); !ok {
			return &CacheInfo{}
		} else if userInfo, ok := cacheObj.(*CacheInfo); ok {
			return userInfo
		}
	}
	return &CacheInfo{}
}

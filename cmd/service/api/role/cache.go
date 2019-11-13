package role

import (
	"context"
	"strconv"
	"time"

	"github.com/fidelfly/gox/cachex/mcache"

	"github.com/fidelfly/fxgos/cmd/utilities/pub"
)

type CacheInfo struct {
	Id          int64
	Code        string
	Description string
	Roles       []int64
}

var myCache = mcache.NewEnsureCache(mcache.DefaultExpiration, 1*time.Hour, roleResolver)

func roleResolver(key string) interface{} {
	id, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return nil
	}
	r, err := Read(context.Background(), id)
	if err != nil {
		return nil
	}

	return &CacheInfo{
		Id:          r.Id,
		Code:        r.Code,
		Description: r.Description,
		Roles:       r.Roles,
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

func GetCache(roleId int64) *CacheInfo {
	if roleId > 0 {
		if cacheObj, ok := myCache.Get(strconv.FormatInt(roleId, 10)); !ok {
			return &CacheInfo{}
		} else if info, ok := cacheObj.(*CacheInfo); ok {
			return info
		}
	}
	return &CacheInfo{}
}

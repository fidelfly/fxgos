package da

import (
	"fmt"
	"os"
	"strings"

	"github.com/fidelfly/gox/cachex/bcache"

	"github.com/fidelfly/fxgos/cmd/utilities/system"
	"github.com/fidelfly/gostool/db"
)

var sgCache *bcache.BuntCache

func encodeUserKey(userId int64) string {
	return fmt.Sprintf("user_%d", userId)
}

func encodeResKey(resType string, resId int64) string {
	return fmt.Sprintf("res_%s_%d", strings.ReplaceAll(resType, ".", "_"), resId)
}

type userSg struct {
	UserId         int64
	SecurityGroups []int64
}

type resourceSg struct {
	ResType        string
	ResId          int64
	SecurityGroups []int64
}

func initSgCache() {
	init := false
	if _, err := os.Stat(system.GetTemporaryPath("sgc")); err != nil {
		init = true
	}
	sgCache, _ = bcache.NewCache(
		system.GetTemporaryPath("sgc"),
		bcache.NewDataset("user_*", bcache.JsonConverter),
		bcache.NewDataset("res_*", bcache.JsonConverter),
	)
	sgCache.CreateJSONIndex("user", "user_*")
	sgCache.CreateJSONIndex("res", "res_*")
	if init {
		_initCache(sgCache)
	}
}

func _initCache(cache *bcache.BuntCache) {
	dbs := db.NewSession()
	defer dbs.Close()
	userSgs := make([]*userSg, 0)
	err := dbs.Find(&userSgs,
		db.SQL("select user_id, JSON_ARRAYAGG(security_group) as security_groups from user_security_group group by user_id"),
	)
	if err != nil {
		return
	}
	if len(userSgs) > 0 {
		for _, usg := range userSgs {
			_ = cache.Set(encodeUserKey(usg.UserId), usg)
		}
	}

	resSgs := make([]*resourceSg, 0)
	err = dbs.Find(&resSgs,
		db.SQL("select res_type, res_id, JSON_ARRAYAGG(security_group) as security_groups from res_security_group group by res_type, res_id"),
	)
	if err != nil {
		return
	}
	if len(resSgs) > 0 {
		for _, rsg := range resSgs {
			_ = cache.Set(encodeResKey(rsg.ResType, rsg.ResId), rsg)
		}
	}
}

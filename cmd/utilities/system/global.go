//todo delete
package system

import (
	"github.com/fidelfly/fxgo"
	"github.com/fidelfly/fxgo/cachex/mcache"
	"github.com/go-xorm/xorm"
)

var UserCache *mcache.MemCache

// nolint:gosec
const (
	TokenPath       = "/fxgos/token"
	ProtectedPrefix = "/fxgos"
	PublicPrefix    = "/public"
)

var DbEngine *xorm.Engine

var TokenKeeper = &fxgo.TokenIssuer{}

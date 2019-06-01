package system

import (
	"github.com/fidelfly/fxgo/cachex"
	"github.com/go-xorm/xorm"
)

var UserCache *cachex.MemCache

const TokenPath = "/fxgos/token"
const ProtectedPrefix = "/fxgos"
const PublicPrefix = "/public"

var DbEngine *xorm.Engine

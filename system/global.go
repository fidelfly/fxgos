package system

import (
	"github.com/fidelfly/fxgo/cachex"
	"github.com/go-xorm/xorm"
)

var UserCache *cachex.MemCache

// nolint:gosec
const (
	TokenPath       = "/fxgos/token"
	ProtectedPrefix = "/fxgos"
	PublicPrefix    = "/public"
)

var DbEngine *xorm.Engine

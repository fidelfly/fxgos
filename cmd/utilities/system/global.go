//todo delete
package system

import (
	"github.com/fidelfly/gox/cachex/mcache"
	"github.com/fidelfly/gox/gosrvx"
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

var TokenKeeper = &gosrvx.TokenIssuer{}

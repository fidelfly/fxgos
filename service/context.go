package service

import (
	"net/http"
	"github.com/lyismydg/fxgos/caches"
)

type contextKeys struct {
	UserInfo string
}

var ContextKeys = contextKeys{
	"CONTEXT_USER_INFO",
}

func GetUserInfo(r *http.Request) *caches.UserInfo {
	user := ContextGet(r, ContextKeys.UserInfo)

	if user !=  nil {
		return user.(*caches.UserInfo)
	}

	return nil
}

package service

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fidelfly/fxgo/httprxr"
	"github.com/fidelfly/fxgo/lockx"

	"github.com/fidelfly/fxgos/system"

	"github.com/fidelfly/fxgos/caches"
)

type contextKeys struct {
	UserInfo string
}

var ContextKeys = contextKeys{
	"CONTEXT_USER_INFO",
}

func GetUserInfo(r *http.Request) *caches.UserInfo {
	user := httprxr.ContextGet(r, ContextKeys.UserInfo)

	if user != nil {
		return user.(*caches.UserInfo)
	}

	return nil
}

//export
func ResourceLockedError(action lockx.Action) httprxr.ResponseMessage {
	var data map[string]interface{}
	if action != nil {
		if userId, err := strconv.ParseInt(action.GetOwnerKey(), 10, 64); err != nil {
			panic(errors.New("owner's key of lock action can not be converted to int64"))
		} else {
			data = make(map[string]interface{})
			user := system.User{
				Id: userId,
			}
			_, err := system.DbEngine.Get(&user)
			if err != nil {
				data["user"] = userId
			} else {
				data["user"] = user.Name
			}
			data["action"] = action.GetCode()
		}

	}

	return httprxr.NewErrorMessage("RESOURCE_LOCKED", "Resource is locked by someone. Please try again later.", data)
}

package router

import (
	"net/http"
	"github.com/lyismydg/fxgos/service"
	"github.com/lyismydg/fxgos/system"
	"github.com/lyismydg/fxgos/auth"
)

const PASSWORD_UNCHANGE  = "PASSWORD_UNCHANGE"
const INVALID_ORG_PASSWORD  = "INVALID_ORG_PASSWORD"

func updatePassword(w http.ResponseWriter, r *http.Request) {
	params := service.GetRequestVars(r, "orgPwd", "newPwd")
	newPwd := params["newPwd"]
	orgPwd := params["orgPwd"]

	if len(newPwd) == 0{
		service.ResponseJSON(w, nil, service.InvalidParamError("newPwd"), http.StatusBadRequest)
		return
	}
	if len(orgPwd) == 0 {
		service.ResponseJSON(w, nil, service.InvalidParamError("orgPwd"), http.StatusBadRequest)
		return
	}

	if newPwd == orgPwd {
		service.ResponseJSON(w, nil, service.NewResponseError(PASSWORD_UNCHANGE, "New password is same as the original password."), http.StatusBadRequest)
		return
	}

	userInfo := service.GetUserInfo(r)

	if userInfo == nil {
		service.ResponseJSON(w, nil, service.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	user := &system.User{
		Id: userInfo.Id,
	}

	if ok, err := system.DbEngine.Get(user); ok {
		if auth.EncodePassword(user.Code, orgPwd) != user.Password {
			service.ResponseJSON(w,nil, service.NewResponseError(INVALID_ORG_PASSWORD, "Original password is wrong"), http.StatusBadRequest)
			return
		}

		user.Password = auth.EncodePassword(user.Code, newPwd)

		if _, err := system.DbEngine.Update(user); err != nil {
			service.ResponseJSON(w, nil, service.ExceptionError(err), http.StatusInternalServerError)
			return
		} else {
			data := make(map[string]interface{})
			data["status"] = "ok"
			service.ResponseJSON(w, nil, data, http.StatusOK)
			return
		}

	} else {
		service.ResponseJSON(w, nil, service.ExceptionError(err), http.StatusInternalServerError)
		return
	}

}

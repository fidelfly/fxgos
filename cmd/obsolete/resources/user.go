package resources

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fidelfly/fxgo/authx"
	"github.com/fidelfly/fxgo/httprxr"
	"github.com/fidelfly/fxgo/logx"

	"github.com/fidelfly/fxgos/cmd/utilities/auth"

	"github.com/fidelfly/fxgos/cmd/utilities/system"
)

type ResourceUser struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Avatar int64  `json:"avatar"`
}
type UserService struct {
}

const duplicateUserCode = "duplicate_user_code"

func (us *UserService) Get(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "userId")
	userID, _ := strconv.ParseInt(params["userId"], 10, 64)
	if userID == 0 {
		user := GetUserInfo(r)
		userID = user.ID
	}
	if userID > 0 {
		user := new(ResourceUser)
		ok, err := system.DbEngine.SQL("select a.id, a.code, a.name, a.avatar from user as a where a.id = ?", userID).Get(user)
		if ok {
			httprxr.ResponseJSON(w, http.StatusOK, user)
			return
		}
		httprxr.ResponseJSON(w, http.StatusNotFound, httprxr.ExceptionMessage(err))
	}
	httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("userID"))
}

func (us *UserService) Post(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "userId", "name", "avatar")
	userID, _ := strconv.ParseInt(params["userId"], 10, 64)

	if userID == 0 {
		userInfo := GetUserInfo(r)
		userID = userInfo.ID
	}

	user := &system.User{
		ID: userID,
	}
	find, err := system.DbEngine.Get(user)
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	}
	if find {
		if len(params["name"]) > 0 {
			user.Name = params["name"]
		}
		if len(params["avatar"]) > 0 {
			user.Avatar, err = strconv.ParseInt(params["avatar"], 10, 64)
			if err != nil {
				httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("avatar"))
				return
			}
		}
		_, err = system.DbEngine.Update(user)
		if err != nil {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
			return
		}

		userRes := ResourceUser{
			ID:     user.ID,
			Code:   user.Code,
			Name:   user.Name,
			Avatar: user.Avatar,
		}

		httprxr.ResponseJSON(w, http.StatusOK, userRes)
		return
	}

	httprxr.ResponseJSON(w, http.StatusNotFound, httprxr.ExceptionMessage(errors.New("record is not found")))

}

func (us *UserService) Register(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	name := r.FormValue("name")
	password := r.FormValue("password")
	if len(code) == 0 {
		httprxr.ResponseJSON(w, http.StatusOK, httprxr.InvalidParamError("code"))
		return
	}
	if len(password) == 0 {
		httprxr.ResponseJSON(w, http.StatusOK, httprxr.InvalidParamError("password"))
		return
	}
	if len(name) == 0 {
		name = code
	}

	session := system.DbEngine.NewSession()
	defer session.Close()
	err := session.Begin()
	defer func() {
		if err != nil {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		}
	}()
	if err != nil {
		return
	}

	user := system.User{Code: code}
	if exist, _ := session.Get(&user); exist {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.NewErrorMessage(duplicateUserCode, "Duplicate code is found!"))
		return
	}

	user.Name = name

	password = auth.EncodePassword(code, password)
	// fmt.Println("Password : " + password)
	user.Password = password
	_, err = session.Insert(&user)
	if err != nil {
		logx.CaptureError(session.Rollback())
		return
	}

	err = session.Commit()
	if err != nil {
		return
	}

	data := httprxr.ResponseData{}
	data["user_id"] = user.ID

	httprxr.ResponseJSON(w, http.StatusOK, data)
}

const PasswordUnchange = "PASSWORD_UNCHANGE"
const InvalidOrgPassword = "INVALID_ORG_PASSWORD"

func (us *UserService) updatePassword(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "orgPwd", "newPwd")
	newPwd := params["newPwd"]
	orgPwd := params["orgPwd"]

	if len(newPwd) == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("newPwd"))
		return
	}
	if len(orgPwd) == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("orgPwd"))
		return
	}

	if newPwd == orgPwd {
		httprxr.ResponseJSON(
			w,
			http.StatusBadRequest,
			httprxr.NewErrorMessage(PasswordUnchange, "New password is same as the original password."),
		)
		return
	}

	userInfo := GetUserInfo(r)

	if userInfo == nil {
		httprxr.ResponseJSON(w, http.StatusUnauthorized, httprxr.NewErrorMessage(authx.UnauthorizedErrorCode, "unauthorized"))
		return
	}

	user := &system.User{
		ID: userInfo.ID,
	}

	if ok, err := system.DbEngine.Get(user); ok {
		if auth.EncodePassword(user.Code, orgPwd) != user.Password {
			httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.NewErrorMessage(InvalidOrgPassword, "Original password is wrong"))
			return
		}

		user.Password = auth.EncodePassword(user.Code, newPwd)

		if _, terr := system.DbEngine.Update(user); terr != nil {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
			return
		}
		data := make(map[string]interface{})
		data["status"] = "ok"
		httprxr.ResponseJSON(w, http.StatusOK, data)
		return

	} else if err != nil {
		logx.Error(err)
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	}
	httprxr.ResponseJSON(w, http.StatusNotFound, nil)

}

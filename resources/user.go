package resources

import (
	"net/http"
	"github.com/lyismydg/fxgos/service"
	"strconv"
	"github.com/lyismydg/fxgos/system"
	"errors"
	"github.com/lyismydg/fxgos/auth"
)

type ResourceUser struct {
	Id int64 `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Avatar int64 `json:"avatar"`
}
type UserService struct {

}

func (us *UserService) Get(w http.ResponseWriter, r *http.Request) {
	params := service.GetRequestVars(r, "userId")
	userId, _ := strconv.ParseInt(params["userId"],10, 64)
	if userId == 0 {
		user := service.GetUserInfo(r)
		userId = user.Id
	}
	if userId > 0 {
		user := new(ResourceUser)
		ok, err := system.DbEngine.SQL("select a.id, a.code, a.name, a.avatar from user as a where a.id = ?", userId).Get(user)
		if ok {
			service.ResponseJSON(w, nil, user, http.StatusOK)
			return
		}
		service.ResponseJSON(w, nil, service.ExceptionError(err), http.StatusNotFound)
	}
	service.ResponseJSON(w, nil, service.InvalidParamError("userId"), http.StatusBadRequest)
}

func (us *UserService) Post(w http.ResponseWriter, r *http.Request) {
	params := service.GetRequestVars(r,"userId", "name", "avatar")
	userId, _ := strconv.ParseInt(params["userId"], 10, 64)

	if userId == 0 {
		userInfo := service.GetUserInfo(r)
		userId = userInfo.Id
	}

	user := &system.User{
		Id: userId,
	}
	find, err := system.DbEngine.Get(user)
	if err != nil {
		service.ResponseJSON(w, nil, service.ExceptionError(err), http.StatusInternalServerError)
		return
	}
	if find {
		if len(params["name"]) > 0 {
			user.Name = params["name"]
		}
		if len(params["avatar"]) > 0 {
			user.Avatar, err = strconv.ParseInt(params["avatar"], 10, 64)
			if err != nil {
				service.ResponseJSON(w, nil, service.InvalidParamError("avatar"), http.StatusBadRequest)
				return
			}
		}
		_, err = system.DbEngine.Update(user)
		if err != nil {
			service.ResponseJSON(w, nil, service.ExceptionError(err), http.StatusInternalServerError)
			return
		}

		userRes := ResourceUser{
			Id: user.Id,
			Code: user.Code,
			Name: user.Name,
			Avatar: user.Avatar,
		}

		service.ResponseJSON(w, nil, userRes, http.StatusOK)
		return
	}

	service.ResponseJSON(w, nil, service.ExceptionError(errors.New("Record Not Found!")), http.StatusNotFound)

}

func (us *UserService) Register(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	name := r.FormValue("name")
	password := r.FormValue("password")
	if len(code) == 0 {
		service.ResponseJSON(w, nil, service.InvalidParamError("code"), http.StatusOK)
		return
	}
	if len(password) == 0 {
		service.ResponseJSON(w, nil, service.InvalidParamError("password"), http.StatusOK)
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
			service.ResponseJSON(w, nil, service.ExceptionError(err), http.StatusInternalServerError)
		}
	}()
	if err != nil {
		return
	}

	password = auth.EncodePassword(code, password)
	//fmt.Println("Password : " + password)
	user := system.User{Code:code, Name:name, Password: password}
	_, err = session.Insert(&user)
	if err != nil {
		session.Rollback()
		return
	}

	err = session.Commit()
	if err != nil {
		return
	}

	data := service.ResponseData{}
	data["user_id"] = user.Id

	service.ResponseJSON(w, nil, data, http.StatusOK)
}


func init() {
	user := new(UserService)
	defineResourceHandlerFunction("get", "/fxgos/user", user.Get)
	defineResourceHandlerFunction("post", "/fxgos/user", user.Post)
	defineResourceHandlerFunction("post", "/public/user", user.Register)
}

package resources

import (
	"net/http"
	"github.com/lyismydg/fxgos/service"
	"strconv"
	"github.com/lyismydg/fxgos/system"
		"github.com/lyismydg/fxgos/auth"
)

type ResourceUser struct {
	Id int64 `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
type UserService struct {

}

func (us *UserService) Get(w http.ResponseWriter, r *http.Request) {
	params := service.GetRequestVars(r, "userId")
	userId, _ := strconv.ParseInt(params["userId"],10, 64)
	if userId > 0 {
		user := new(ResourceUser)
		ok, err := system.DbEngine.SQL("select a.id, a.code, a.name from user as a where a.id = ?", userId).Get(user)
		if ok {
			service.ResponseJSON(w, nil, user, http.StatusOK)
			return
		}
		service.ResponseJSON(w, nil, service.ExceptionError(err), http.StatusNotFound)
	}
	service.ResponseJSON(w, nil, service.InvalidParamError("userId"), http.StatusBadRequest)
}

func (us *UserService) Post(w http.ResponseWriter, r *http.Request) {
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
	path := service.GetProtectedPath("user")
	defineResourceHandlerFunction("get", path, user.Get)
	defineResourceHandlerFunction("post", path, user.Post)
}

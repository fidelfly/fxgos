package resources

import (
	"net/http"
	"github.com/lyismydg/fxgos/service"
	"strconv"
	"github.com/lyismydg/fxgos/system"
)

type ResourceUser struct {
	Id int64 `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	TenantId int64 `json:"tenant_id"`
	TenantName string `json:"tenant_name"`
	TenantCode string `json:"tenant_code"`
}
type UserService struct {

}

func (us *UserService) Get(w http.ResponseWriter, r *http.Request) {
	params := service.GetRequestVars(r, "userId")
	userId, _ := strconv.ParseInt(params["userId"],10, 64)
	if userId > 0 {
		user := new(ResourceUser)
		ok, err := system.DbEngine.SQL("select a.id, a.code, a.name, a.tenant_id, b.code as tenant_code, b.name as tenant_name from user as a, tenant as b where a.id = ? and a.tenant_id = b.id", userId).Get(user)
		if ok {
			service.ResponseJSON(w, nil, user, http.StatusOK)
			return
		}
		service.ResponseJSON(w, nil, service.ExceptionError(err), http.StatusNotFound)
	}
	service.ResponseJSON(w, nil, service.InvalidParamError("userId"), http.StatusBadRequest)
}


func init() {
	user := new(UserService)
	defineResourceHandlerFunction("get", "/admin/user", user.Get)
}

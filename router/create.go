package router

import (
	"net/http"
	"github.com/lyismydg/fxgos/system"
	"github.com/lyismydg/fxgos/service"
	"fmt"
	"github.com/lyismydg/fxgos/auth"
)

func createTeant(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	name := r.FormValue("name")
	password := r.FormValue("password")
	if len(code) == 0 {
		service.ResponseJSON(w, nil, service.InvalidParamError("code"), http.StatusOK)
		return
	}
	tenant := system.Tenant{Code: code, Name: name}
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

	_, err = session.Insert(&tenant)
	if err != nil {
		session.Rollback()
		return
	}
	if len(password) == 0 {
		password = fmt.Sprintf("%s123456", code)
	}
	password = auth.EncodePassword(code, password)
	fmt.Println("Password : " + password)
	user := system.User{Code:code, Name:name, TenantId: tenant.Id, Password: password}
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
	data["tenant_id"] = tenant.Id
	data["user_id"] = user.Id

	service.ResponseJSON(w, nil, data, http.StatusOK)
}
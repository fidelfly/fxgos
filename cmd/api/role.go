package api

import (
	"net/http"
	"strconv"

	"github.com/fidelfly/gox/httprxr"
	"github.com/fidelfly/gox/logx"
	"github.com/fidelfly/gox/routex"

	"github.com/fidelfly/fxgos/cmd/service/iam"
	"github.com/fidelfly/fxgos/cmd/service/iam/iamx"
	"github.com/fidelfly/fxgos/cmd/service/role"
	"github.com/fidelfly/fxgos/cmd/service/role/res"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/gostool/dbo"
)

func RoleRoute(router *routex.Router) {
	rr := router.PathPrefix("/role").Subrouter()
	rr.Restricted(true)

	resource := "role"
	accessRight := iam.NewAccessItem(iamx.ResourceFunction, resource, iamx.ActionAccess)
	createRight := iam.NewAccessItem(iamx.ResourceFunction, resource, iamx.ActionCreate)
	updateRight := iam.NewAccessItem(iamx.ResourceFunction, resource, iamx.ActionUpdate)
	deleteRight := iam.NewAccessItem(iamx.ResourceFunction, resource, iamx.ActionDelete)
	rr.Methods(http.MethodGet).Path("/policy").HandlerFunc(listPolicy).ApplyProps(iamProps(accessRight))
	rr.Methods(http.MethodGet).Path("/acl").HandlerFunc(listRoleAcl).ApplyProps(iamProps(accessRight))
	rr.Methods(http.MethodGet).Path("/list").HandlerFunc(listRole).ApplyProps(iamProps(accessRight))
	rr.Methods(http.MethodGet).Path("/{id}").HandlerFunc(getRole).ApplyProps(iamProps(accessRight))
	rr.Methods(http.MethodGet).HandlerFunc(getRole).ApplyProps(iamProps(accessRight))
	rr.Methods(http.MethodPost).HandlerFunc(postRole).ApplyProps(iamProps(createRight))
	rr.Methods(http.MethodPut).HandlerFunc(putRole).ApplyProps(iamProps(updateRight))
	rr.Methods(http.MethodDelete).Path("/{id}").HandlerFunc(deleteRole).ApplyProps(iamProps(deleteRight))
	rr.Methods(http.MethodDelete).HandlerFunc(deleteRole).ApplyProps(iamProps(deleteRight))
}

func listRoleAcl(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "role_id")
	roleId, _ := strconv.ParseInt(params["role_id"], 10, 64)
	if roleId == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("role_id"))
		return
	}
	iamPolicys := iam.ListResourceAclByRole(r.Context(), roleId, iamx.ResourceFunction)
	httprxr.ResponseJSON(w, http.StatusOK, iamPolicys)
}

func listPolicy(w http.ResponseWriter, r *http.Request) {
	iamPolicys := iam.ListResourceAclByRole(r.Context(), 0, iamx.ResourceFunction)
	httprxr.ResponseJSON(w, http.StatusOK, iamPolicys)
}

func listRole(w http.ResponseWriter, r *http.Request) {
	params := httprxr.ParseRequestVars(r, "results", "page", "sortField", "sortOrder")
	if listInfo, err := NewList(params); err != nil {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.ExceptionMessage(err))
		return
	} else {
		rsp, count, err := role.List(r.Context(), listInfo)
		if err != nil {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
			return
		}
		data := make(map[string]interface{})
		data["count"] = count
		data["data"] = rsp
		httprxr.ResponseJSON(w, http.StatusOK, data)
	}
}

func getRole(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "id")
	roleId, _ := strconv.ParseInt(params["id"], 10, 64)
	if roleId == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}
	if rsp, err := role.Read(r.Context(), roleId); err != nil {
		if err == syserr.ErrNotFound {
			httprxr.ResponseJSON(w, http.StatusNotFound, nil)
		} else {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		}
	} else {
		httprxr.ResponseJSON(w, http.StatusOK, rsp)
	}
}

type RoleInput struct {
	res.Role
	IamPolicys []*iam.ResourceACL `json:"iamPolicys"` //todo adjust name of variable
}

func postRole(w http.ResponseWriter, r *http.Request) {
	task, ok := registerTask(TaskMenuRole, TaskOperationCreate, r, w)
	if ok {
		defer task.Done()
	} else {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, []byte("Cannot register task"))
		return
	}
	roleInput := &RoleInput{}
	if err := httprxr.GetJSONRequestData(r, roleInput); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	}

	if len(roleInput.Code) == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("code"))
		return
	}

	task.StartTrail()
	task.SetField("RoleCode", roleInput.Code)
	ctx, dbs := dbo.WithDBSession(r.Context())
	defer dbs.Close()
	if err := dbs.Begin(); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(syserr.DatabaseErr(err)))
	}
	if rsp, err := role.Create(ctx, role.Form{
		Code:        roleInput.Code,
		Roles:       roleInput.Roles,
		Description: roleInput.Description,
	}); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		task.LogTrailDone(err)
		logx.CaptureError(dbs.Rollback())
		return
	} else {
		task.LogTrailDone(nil)
		if err := iam.UpdatePolicyByRole(ctx, rsp.Id, rsp.Roles, roleInput.IamPolicys); err != nil {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(syserr.DatabaseErr(err)))
			logx.CaptureError(dbs.Rollback())
			return
		}
		if err := dbs.Commit(); err != nil {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
			return
		}
		httprxr.ResponseJSON(w, http.StatusOK, rsp)
	}
}

func putRole(w http.ResponseWriter, r *http.Request) {
	task, ok := registerTask(TaskMenuRole, TaskOperationModify, r, w)
	if ok {
		defer task.Done()
	} else {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, []byte("Cannot register task"))
		return
	}
	roleInput := &RoleInput{}
	if err := httprxr.GetJSONRequestData(r, roleInput); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	}
	roleData := &res.Role{
		Id:          roleInput.Id,
		Code:        roleInput.Code,
		Description: roleInput.Description,
		Roles:       roleInput.Roles,
	}
	if roleData.Id == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}

	task.StartTrail()
	task.SetField("RoleCode", roleData.Code)
	ctx, dbs := dbo.WithDBSession(r.Context())
	defer dbs.Close()
	if err := dbs.Begin(); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(syserr.DatabaseErr(err)))
	}
	if err := role.Update(ctx, dbo.UpdateInfo{
		Id:   roleData.Id,
		Cols: []string{"description", "roles"},
		Data: roleData,
	}); err != nil {
		if err == syserr.ErrNotFound {
			httprxr.ResponseJSON(w, http.StatusNotFound, nil)
		} else {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		}
		logx.CaptureError(dbs.Rollback())
		task.LogTrailDone(err)
	} else {
		task.LogTrailDone(nil)
		if err := iam.UpdatePolicyByRole(ctx, roleData.Id, roleData.Roles, roleInput.IamPolicys); err != nil {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
			logx.CaptureError(dbs.Rollback())
			return
		}
		if err := dbs.Commit(); err != nil {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
			return
		}
		httprxr.ResponseJSON(w, http.StatusOK, nil)
	}
}

func deleteRole(w http.ResponseWriter, r *http.Request) {
	task, ok := registerTask(TaskMenuRole, TaskOperationDelete, r, w)
	if ok {
		defer task.Done()
	} else {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, []byte("Cannot register task"))
		return
	}
	params := httprxr.GetRequestVars(r, "id")
	roleId, _ := strconv.ParseInt(params["id"], 10, 64)
	if roleId == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}

	task.StartTrail()
	if rsp, err := role.Read(r.Context(), roleId); err != nil {
		if err == syserr.ErrNotFound {
			httprxr.ResponseJSON(w, http.StatusNotFound, nil)
		} else {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		}
		return
	} else {
		task.SetField("RoleCode", rsp.Code)
	}

	if err := role.Delete(r.Context(), roleId); err != nil {
		if err == syserr.ErrNotFound {
			httprxr.ResponseJSON(w, http.StatusNotFound, nil)
		} else {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		}
		task.LogTrailDone(err)
	} else {
		task.LogTrailDone(nil)
		if err := iam.DeleteRolePolicy(r.Context(), roleId); err != nil {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		} else {
			httprxr.ResponseJSON(w, http.StatusOK, nil)
		}
	}
}

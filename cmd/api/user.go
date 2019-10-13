package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/fidelfly/gox/authx"
	"github.com/fidelfly/gox/httprxr"
	"github.com/fidelfly/gox/logx"
	"github.com/fidelfly/gox/routex"

	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/pkg/mail"
	"github.com/fidelfly/fxgos/cmd/service/iam"
	"github.com/fidelfly/fxgos/cmd/service/iam/iamx"
	"github.com/fidelfly/fxgos/cmd/service/otk"
	"github.com/fidelfly/fxgos/cmd/service/user"
	"github.com/fidelfly/fxgos/cmd/service/user/res"
	"github.com/fidelfly/fxgos/cmd/utilities/auth"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/fxgos/cmd/utilities/system"
)

func UserRoute(router *routex.Router) {
	rr := router.PathPrefix("/user").Subrouter()
	rr.Restricted(true)

	resource := "users"
	accessRight := iam.NewAccessItem(iamx.ResourceFunction, resource, iamx.ActionAccess)
	createRight := iam.NewAccessItem(iamx.ResourceFunction, resource, iamx.ActionCreate)
	updateRight := iam.NewAccessItem(iamx.ResourceFunction, resource, iamx.ActionUpdate)
	deleteRight := iam.NewAccessItem(iamx.ResourceFunction, resource, iamx.ActionDelete)

	rr.Methods(http.MethodPost).Path("/password").HandlerFunc(updatePassword)
	rr.Methods(http.MethodPost).Path("/activeEmail").HandlerFunc(sendActivateMail).Restricted(false)
	rr.Methods(http.MethodPost).Path("/activate").HandlerFunc(activateUser).Restricted(false)
	rr.Methods(http.MethodPost).Path("/resetPwdEmail").HandlerFunc(resetPwdEmail).Restricted(false)
	rr.Methods(http.MethodPost).Path("/resetPwd").HandlerFunc(resetPwd).Restricted(false)

	AttachAccessPremises(rr.Methods(http.MethodGet).Path("/list").HandlerFunc(listUser), accessRight)
	AttachAccessPremises(rr.Methods(http.MethodPost).Path("/disable/{id}").HandlerFunc(disableUser), updateRight)
	AttachAccessPremises(rr.Methods(http.MethodPost).Path("/enable/{id}").HandlerFunc(enableUser), updateRight)
	AttachAccessPremises(rr.Methods(http.MethodPost).Path("/activateEmail/{id}").HandlerFunc(activateEmail), updateRight)
	rr.Methods(http.MethodGet).Path("/{id}").HandlerFunc(getUser)
	rr.Methods(http.MethodGet).Path("/acl/{id}").HandlerFunc(listUserAcl)
	rr.Methods(http.MethodGet).Path("/acl").HandlerFunc(listUserAcl)
	rr.Methods(http.MethodGet).HandlerFunc(getUser)
	rr.Methods(http.MethodPut).HandlerFunc(putUser)
	AttachAccessPremises(rr.Methods(http.MethodPost).HandlerFunc(postUser), createRight)
	AttachAccessPremises(rr.Methods(http.MethodDelete).Path("/{id}").HandlerFunc(deleteUser), deleteRight)

	router.Path("/logout").HandlerFunc(auth.TokenIssuer.AuthorizeDisposeHandlerFunc).Restricted(true)
}

func updatePassword(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "orgPwd", "newPwd")
	if !CheckEmptyParam(w, r, params, "orgPwd", "newPwd") {
		return
	}
	orgPwd := params["orgPwd"]
	newPwd := params["newPwd"]
	if orgPwd == newPwd {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.NewErrorMessage("PASSWORD_UNCHANGE", "New password is same as the original password."))
		return
	}

	currentUser := GetUserInfo(r)
	if currentUser == nil {
		httprxr.ResponseJSON(w, http.StatusUnauthorized, httprxr.NewErrorMessage(authx.UnauthorizedErrorCode, "invalid access"))
		return
	}

	if _, err := user.Validate(r.Context(), user.ValidateInput{
		Id:       currentUser.Id,
		Password: orgPwd,
	}); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.NewErrorMessage("INVALID_ORG_PASSWORD", "Original password is wrong"))
		return
	}

	if id, err := user.Update(r.Context(), user.UpdateInput{
		UpdateInfo: db.UpdateInfo{Id: currentUser.Id, Cols: []string{"password"}},
		Data:       &res.User{Id: currentUser.Id, Password: newPwd},
	}); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
	} else if id == 0 {
		httprxr.ResponseJSON(w, http.StatusNotFound, nil)
	} else {
		httprxr.ResponseJSON(w, http.StatusOK, nil)
	}
}

func sendActivateMail(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "id", "username", "password")
	var userData *res.User
	var err error
	if len(params["id"]) > 0 {
		userID, _ := strconv.ParseInt(params["id"], 10, 64)
		if userID == 0 {
			httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
			return
		}
		userData, err = user.Read(r.Context(), userID)
	} else {
		if !CheckEmptyParam(w, r, params, "username", "password") {
			return
		}

		userData, err = user.Validate(r.Context(), user.ValidateInput{Code: params["username"], Password: params["password"]})
	}
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	}
	if userData == nil {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}

	if userData.Status != user.StatusDeactivated {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}

	if key, err := otk.NewOtk("USER_ACTIVE", strconv.FormatInt(userData.Id, 10), 24*time.Hour, "Active User", otk.NewResourceKey(userData.Id)); err == nil {
		data := make(map[string]string)
		data["name"] = userData.Name
		if len(params["id"]) > 0 {
			data["activelink"] = system.Runtime.Domain + "/sys/activate?otk=" + url.QueryEscape(key)
		} else {
			data["activelink"] = system.Runtime.Domain + "/api/user/activate?otk=" + url.QueryEscape(key)
		}

		if err := mail.SendMail(mail.CreateMessage(
			mail.TemplateMessage("", "activemail.tpl", data),
			mail.Subject("Activate your account"),
			mail.To(userData.Email))); err != nil {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		} else {
			httprxr.ResponseJSON(w, http.StatusOK, map[string]interface{}{"email": userData.Email})
		}
	} else {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
	}
}

func activateUser(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "otk", "pwd")
	if !CheckEmptyParam(w, r, params, "otk") {
		return
	}

	otkData, err := otk.Validate(params["otk"])
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	} else if otkData.Id == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest,
			httprxr.NewErrorMessage("INVALID_RESETPWD_LINK", "Invalid Reset Pwd link."))
		return
	}

	userId, _ := strconv.ParseInt(otkData.TypeId, 10, 64)
	updateCols := []string{"status"}
	data := &res.User{
		Id:     userId,
		Status: user.StatusValid,
	}
	if len(params["pwd"]) != 0 {
		updateCols = append(updateCols, "password")
		data.Password = params["pwd"]
	}

	if tarId, err := user.Update(r.Context(), user.UpdateInput{
		UpdateInfo: db.UpdateInfo{Id: data.Id, Cols: updateCols},
		Data:       data,
	}); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	} else if tarId == 0 {
		httprxr.ResponseJSON(w, http.StatusNotFound, nil)
		return
	}

	if err := otk.Consume(otkData.Id); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
	} else {
		httprxr.ResponseJSON(w, http.StatusOK, nil)
	}
}

func resetPwdEmail(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "userEmail")
	if !CheckEmptyParam(w, r, params, "userEmail") {
		return
	}
	userData, err := user.ReadByEmail(r.Context(), params["userEmail"])
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
	} else if userData.Id == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.NewErrorMessage("USER_NOT_EXIST", "Can't find this user."))
	}

	if key, err := otk.NewOtk("USER_RESETPWD", strconv.FormatInt(userData.Id, 10), 24*time.Hour, "Reset Password", otk.NewResourceKey(userData.Id)); err == nil {
		data := make(map[string]string)
		data["name"] = userData.Name
		data["resetPwdLink"] = system.Runtime.Domain + "/resetPwd?otk=" + url.QueryEscape(key)
		if err := mail.SendMail(mail.CreateMessage(
			mail.TemplateMessage("", "resetPwdMail.tpl", data),
			mail.Subject("[fxgos] Please reset your password"),
			mail.To(userData.Email))); err != nil {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		}
	} else {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
	}
}

func resetPwd(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "otk", "pwd")
	if !CheckEmptyParam(w, r, params, "otk", "pwd") {
		return
	}

	otkData, err := otk.Validate(params["otk"])
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	} else if otkData.Id == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.NewErrorMessage("INVALID_RESETPWD_LINK", "Invalid Reset Pwd link."))
		return
	}

	userId, _ := strconv.ParseInt(otkData.TypeId, 10, 64)
	if tarId, err := user.Update(r.Context(), user.UpdateInput{
		UpdateInfo: db.UpdateInfo{
			Id:   userId,
			Cols: []string{"password"},
		},
		Data: &res.User{
			Id:       userId,
			Password: params["pwd"],
		},
	}); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	} else if tarId == 0 {
		httprxr.ResponseJSON(w, http.StatusNotFound, nil)
		return
	}

	if err := otk.Consume(otkData.Id); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
	} else {
		httprxr.ResponseJSON(w, http.StatusOK, nil)
	}
}

func listUser(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "results", "page", "sortField", "sortOrder", "statusType", "includedDel")
	statusType := params["statusType"]
	includedDel, _ := strconv.ParseBool(params["includedDel"])
	cond := ""
	if len(statusType) > 0 && statusType != "--" {
		cond = fmt.Sprintf("status = '%s'", statusType)
	} else if !includedDel {
		cond = fmt.Sprintf("status != %d", user.StatusDeleted)
	}
	list, count, err := user.List(r.Context(), NewListInfo(params, cond))
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	}
	data := make(map[string]interface{})
	data["count"] = count
	data["data"] = list
	httprxr.ResponseJSON(w, http.StatusOK, data)
}

func disableUser(w http.ResponseWriter, r *http.Request) {
	task, ok := registerTask(TaskMenuUsers, TaskOperationDisable, r, w)
	if ok {
		defer task.Done()
	} else {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, []byte("Cannot register task"))
		return
	}
	params := httprxr.GetRequestVars(r, "id")
	userID, _ := strconv.ParseInt(params["id"], 10, 64)
	if userID == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}
	currentUser := GetUserInfo(r)
	if currentUser != nil && userID == currentUser.Id {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}

	task.StartTrail()
	task.SetField("UserName", user.GetCache(userID).Name)
	if tarId, err := user.Update(r.Context(), user.UpdateInput{
		UpdateInfo: db.UpdateInfo{
			Id:   userID,
			Cols: []string{"status"},
		},
		Data: &res.User{
			Id:     userID,
			Status: user.StatusInvalid,
		},
	}); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		task.LogTrailDone(err)
	} else if tarId == 0 {
		httprxr.ResponseJSON(w, http.StatusNotFound, nil)
	} else {
		httprxr.ResponseJSON(w, http.StatusOK, nil)
		task.LogTrailDone(nil)
	}
}

func enableUser(w http.ResponseWriter, r *http.Request) {
	task, ok := registerTask(TaskMenuUsers, TaskOperationEnable, r, w)
	if ok {
		defer task.Done()
	} else {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, []byte("Cannot register task"))
		return
	}
	params := httprxr.GetRequestVars(r, "id")
	userID, _ := strconv.ParseInt(params["id"], 10, 64)
	if userID == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}
	currentUser := GetUserInfo(r)
	if currentUser != nil && userID == currentUser.Id {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}
	userData, err := user.Read(r.Context(), userID)
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	} else if userData.Id == 0 {
		httprxr.ResponseJSON(w, http.StatusNotFound, nil)
		return
	}

	if userData.Status != user.StatusInvalid {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}

	task.StartTrail()
	task.SetField("UserName", userData.Name)
	if tarId, err := user.Update(r.Context(), user.UpdateInput{
		UpdateInfo: db.UpdateInfo{
			Id:   userID,
			Cols: []string{"status"},
		},
		Data: &res.User{
			Id:     userID,
			Status: user.StatusDeactivated,
		},
	}); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		task.LogTrailDone(err)
		return
	} else if tarId == 0 {
		httprxr.ResponseJSON(w, http.StatusNotFound, nil)
		return
	} else {
		httprxr.ResponseJSON(w, http.StatusOK, nil)
		task.LogTrailDone(nil)
	}

	if key, err := otk.NewOtk("USER_ACTIVE", strconv.FormatInt(userData.Id, 10), 24*time.Hour, "Active User", otk.NewResourceKey(userData.Id)); err == nil {
		data := make(map[string]string)
		data["name"] = userData.Name
		data["activelink"] = system.Runtime.Domain + "/sys/activate?otk=" + url.QueryEscape(key)
		logx.CaptureError(mail.SendMail(mail.CreateMessage(
			mail.TemplateMessage("", "activemail.tpl", data),
			mail.Subject("Activate your account"),
			mail.To(userData.Email),
		)))
	}
}

func activateEmail(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "id", "username", "password")
	var userData *res.User
	var err error
	if len(params["id"]) > 0 {
		userID, _ := strconv.ParseInt(params["id"], 10, 64)
		if userID == 0 {
			httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
			return
		}
		userData, err = user.Read(r.Context(), userID)
	} else {
		if !CheckEmptyParam(w, r, params, "username", "password") {
			return
		}
		userData, err = user.Validate(r.Context(), user.ValidateInput{
			Code:     params["username"],
			Password: params["password"],
		})
	}
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	}
	if userData.Id == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}

	if key, err := otk.NewOtk("USER_ACTIVE", strconv.FormatInt(userData.Id, 10), 24*time.Hour, "Active User", otk.NewResourceKey(userData.Id)); err == nil {
		data := make(map[string]string)
		data["name"] = userData.Name
		if len(params["id"]) > 0 {
			data["activelink"] = system.Runtime.Domain + "/sys/activate?otk=" + url.QueryEscape(key)
		} else {
			data["activelink"] = system.Runtime.Domain + "/api/user/activate?otk=" + url.QueryEscape(key)
		}
		if err := mail.SendMail(mail.CreateMessage(
			mail.TemplateMessage("", "activemail.tpl", data),
			mail.Subject("Activate your account"),
			mail.To(userData.Email),
		)); err != nil {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		} else {
			httprxr.ResponseJSON(w, http.StatusOK, map[string]interface{}{"email": userData.Email})
		}
	} else {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
	}
}

func listUserAcl(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "id")
	userID, _ := strconv.ParseInt(params["id"], 10, 64)
	if userID == 0 {
		currentUser := GetUserInfo(r)
		if currentUser != nil {
			userID = currentUser.Id
		}
	}
	if userID == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}

	accessItems := iam.ListResourceAclByUser(r.Context(), userID, iamx.ResourceFunction)

	httprxr.ResponseJSON(w, http.StatusOK, accessItems)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "id")
	userID, _ := strconv.ParseInt(params["id"], 10, 64)
	if userID == 0 {
		currentUser := GetUserInfo(r)
		if currentUser != nil {
			userID = currentUser.Id
		}
	}
	if userID == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}
	if userData, err := user.Read(r.Context(), userID); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
	} else if userData.Id == 0 {
		httprxr.ResponseJSON(w, http.StatusNotFound, nil)
	} else {
		userData.Password = "" //remove password
		httprxr.ResponseJSON(w, http.StatusOK, userData)
	}
}

func putUser(w http.ResponseWriter, r *http.Request) {
	task, ok := registerTask(TaskMenuUsers, TaskOperationModify, r, w)
	if ok {
		defer task.Done()
	} else {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, []byte("Cannot register task"))
		return
	}
	params := httprxr.ParseRequestVars(r, "id", "name", "avatar", "region", "dept", "tel", "roles")
	userID, err := params.GetInt("id")
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}
	avatar, err := params.GetInt("avatar")
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("avatar"))
		return
	}
	roles, err := params.GetInts("roles")
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("roles"))
		return
	}
	if userID == 0 {
		currentUser := GetUserInfo(r)
		if currentUser != nil {
			userID = currentUser.Id
		}
	}
	if userID == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}

	task.StartTrail()
	task.SetField("UserName", user.GetCache(userID).Name)
	if tarId, err := user.Update(r.Context(), user.UpdateInput{
		UpdateInfo: db.UpdateInfo{
			Id:   userID,
			Cols: []string{"name", "avatar", "region", "dept", "tel", "roles"},
		},
		Data: &res.User{
			Id:     userID,
			Name:   params.GetString("name"),
			Avatar: avatar,
			Region: params.GetString("region"),
			Dept:   params.GetString("dept"),
			Tel:    params.GetString("tel"),
			Roles:  roles,
		},
	}); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		task.LogTrailDone(err)
	} else if tarId == 0 {
		httprxr.ResponseJSON(w, http.StatusNotFound, nil)
	} else {
		httprxr.ResponseJSON(w, http.StatusOK, nil)
		task.LogTrailDone(nil)
	}
}

func postUser(w http.ResponseWriter, r *http.Request) {
	task, ok := registerTask(TaskMenuUsers, TaskOperationCreate, r, w)
	if ok {
		defer task.Done()
	} else {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, []byte("Cannot register task"))
		return
	}
	params := httprxr.ParseRequestVars(r, "code", "name", "email", "password", "region", "dept", "tel", "roles")
	if !CheckEmptyVar(w, r, params, "code") {
		return
	}
	roles, err := params.GetInts("roles")
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("roles"))
		return
	}

	userData := &res.User{
		Code:     params.GetString("code"),
		Name:     params.GetString("name"),
		Email:    params.GetString("email"),
		Password: params.GetString("password"),
		Region:   params.GetString("region"),
		Dept:     params.GetString("dept"),
		Tel:      params.GetString("tel"),
		Roles:    roles,
	}

	task.StartTrail()
	task.SetField("UserName", userData.Name)
	if userId, err := user.Create(r.Context(), userData); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		task.LogTrailDone(err)
		return
	} else {
		userData.Id = userId
		httprxr.ResponseJSON(w, http.StatusOK, userId)
		task.LogTrailDone(nil)
	}

	if key, err := otk.NewOtk("USER_ACTIVE", strconv.FormatInt(userData.Id, 10), 24*time.Hour, "Active User", otk.NewResourceKey(userData.Id)); err == nil {
		data := make(map[string]string)
		data["name"] = userData.Name
		data["activelink"] = system.Runtime.Domain + "/sys/activate?otk=" + url.QueryEscape(key)
		logx.CaptureError(mail.SendMail(mail.CreateMessage(
			mail.TemplateMessage("", "activemail.tpl", data),
			mail.Subject("Activate your account"),
			mail.To(userData.Email),
		)))
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	task, ok := registerTask(TaskMenuUsers, TaskOperationDelete, r, w)
	if ok {
		defer task.Done()
	} else {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, []byte("Cannot register task"))
		return
	}
	params := httprxr.GetRequestVars(r, "id")
	userID, _ := strconv.ParseInt(params["id"], 10, 64)
	if userID == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}
	currentUser := GetUserInfo(r)
	if currentUser != nil && userID == currentUser.Id {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}

	task.StartTrail()
	task.SetField("UserName", user.GetCache(userID).Name)
	if err := user.Delete(r.Context(), userID); err != nil {
		if err == syserr.ErrNotFound {
			httprxr.ResponseJSON(w, http.StatusNotFound, nil)
		} else {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		}
		task.LogTrailDone(err)
	} else {
		httprxr.ResponseJSON(w, http.StatusOK, nil)
		task.LogTrailDone(nil)
	}
}

package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fidelfly/gox/httprxr"
	"github.com/fidelfly/gox/pkg/strh"
	"github.com/fidelfly/gox/routex"

	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
)

func QueryRoute(router *routex.Router) {
	rr := router.PathPrefix("/query").Subrouter()
	rr.Restricted(true)
	rr.Methods(http.MethodGet).Path("/code").HandlerFunc(queryCode)
	rr.Methods(http.MethodGet).Path("/user").HandlerFunc(queryUser)
	rr.Methods(http.MethodGet).Path("/field").HandlerFunc(queryField)
}

func queryCode(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "resource", "code")
	if !CheckEmptyParam(w, r, params, "resource", "code") {
		return
	}
	if exist, err := queryExist(r.Context(), strh.UnderscoreString(params["resource"]), fmt.Sprintf("code = '%s' ", params["code"])); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
	} else {
		httprxr.ResponseJSON(w, http.StatusOK, exist)
	}
}

func queryExist(ctx context.Context, table string, cond string) (bool, error) {
	if len(table) == 0 {
		return false, syserr.ErrInvalidParam
	}
	if exist, err := db.Engine.Table(table).Where(cond).Exist(); err != nil {
		return false, err
	} else {
		return exist, nil
	}
}

func queryUser(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "field", "value")
	if !CheckEmptyParam(w, r, params, "field", "value") {
		return
	}
	if exist, err := queryExist(r.Context(), "user", fmt.Sprintf("%s = '%s' and status != -2 ", params["field"], params["value"])); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
	} else {
		httprxr.ResponseJSON(w, http.StatusOK, exist)
	}
}

func queryField(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "resource", "field", "value", "cond")
	if !CheckEmptyParam(w, r, params, "resource", "field", "value", "cond") {
		return
	}
	if rsp, err := queryExist(r.Context(), strh.UnderscoreString(params["resource"]), fmt.Sprintf("%s = '%s' and %s ", params["field"], params["value"], params["cond"])); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
	} else {
		httprxr.ResponseJSON(w, http.StatusOK, rsp)
	}
}

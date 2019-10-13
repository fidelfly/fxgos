package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fidelfly/gox/httprxr"
	"github.com/fidelfly/gox/pkg/jmap"
	"github.com/fidelfly/gox/routex"

	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/iam"
	"github.com/fidelfly/fxgos/cmd/service/user"
)

const Token = "/api/token"

func CheckEmptyParam(w http.ResponseWriter, r *http.Request, params map[string]string, keys ...string) bool {
	if len(keys) == 0 {
		return true
	}
	for _, key := range keys {
		if len(params[key]) == 0 {
			httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError(key))
			return false
		}
	}
	return true
}

func CheckEmptyVar(w http.ResponseWriter, r *http.Request, vars httprxr.RequestVar, keys ...string) bool {
	if len(keys) == 0 {
		return true
	}
	for _, key := range keys {
		if len(vars.GetString(key)) == 0 {
			httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError(key))
			return false
		}
	}
	return true
}

func NewListInfo(params map[string]string, cond ...string) db.ListInfo {
	req := db.ListInfo{}
	req.Results, _ = strconv.ParseInt(params["results"], 10, 64)
	req.Page, _ = strconv.ParseInt(params["page"], 10, 64)
	req.SortField = params["sortField"]
	req.SortOrder = params["sortOrder"]
	if len(cond) > 0 {
		req.Cond = cond[0]
	}
	return req
}

func AttachAccessPremises(route *routex.Route, premises ...iam.AccessItem) {
	route.SetProps(AccessConfigKey, iam.AccessPremise(premises))
}

func GetPath(r *http.Request) []string {
	//return strings.SplitAfter()
	return nil
}

func userDecorator(idField string, fields ...string) jmap.Decorator {
	return func(path string, maps jmap.JSONMap) {
		if idVal, ok := maps[idField]; ok {
			if id, ok := idVal.(int64); ok {
				if info := user.GetCache(id); info != nil {
					if len(fields) == 0 {
						fields = []string{fmt.Sprintf("%s_name", idField)}
					}
					for _, field := range fields {
						maps[field] = info.Name
					}
				}
			}
		}
	}
}

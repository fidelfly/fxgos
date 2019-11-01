package api

import (
	"fmt"
	"net/http"

	"github.com/fidelfly/gox/httprxr"
	"github.com/fidelfly/gox/pkg/jmap"
	"github.com/fidelfly/gox/routex"

	"github.com/fidelfly/fxgos/cmd/service/iam"
	"github.com/fidelfly/fxgos/cmd/service/user"
	"github.com/fidelfly/gostool/dbo"
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

func ExtractListInfo(param httprxr.RequestVar) (*dbo.ListInfo, error) {
	if !param.Exist("results") {
		return nil, nil
	}
	info := &dbo.ListInfo{Page: 1}
	if result, err := param.GetInt("results"); err != nil {
		return nil, err
	} else {
		info.Results = int(result)
	}
	if param.Exist("page") {
		if page, err := param.GetInt("page"); err != nil {
			return nil, err
		} else {
			info.Page = int(page)
		}
	}

	info.SortField = param.GetString("sortField")
	info.SortOrder = param.GetString("sortOrder")

	return info, nil
}

func iamProps(premises ...iam.AccessItem) routex.PropSetter {
	return func(props *routex.RouteProps) {
		_ = props.Set(AccessConfigKey, iam.AccessPremise(premises))
	}
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

func DeferResponse(w http.ResponseWriter, err error) {
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	}
}

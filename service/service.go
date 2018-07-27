package service

import (
"encoding/json"
"net/http"

"github.com/gorilla/mux"
"context"
"github.com/sirupsen/logrus"
	"github.com/lyismydg/fxgos/system"
	"gopkg.in/oauth2.v3"
	"strings"
)

type ResponseData map[string]interface{}

type RespSetting struct {
	ContentType string
	Cache       bool
	Header      http.Header
}

type ResponseError struct {
	Code string `json:"errorCode"`
	Message string `json:"errorMessage"`
}

var DefaultResponse = &RespSetting{
	ContentType: "application/json;charset=UTF-8",
	Cache:       false,
}

func ResponseJSON(w http.ResponseWriter, setting *RespSetting, data interface{}, statusCode ...int) (err error) {
	if setting == nil {
		setting = DefaultResponse
	}

	if len(setting.ContentType) == 0 {
		setting.ContentType = "application/json;charset=UTF-8"
	}

	w.Header().Set("Content-Type", setting.ContentType)
	if !setting.Cache {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("pragma", "no-cache")
	}

	if setting.Header != nil {
		for key := range setting.Header {
			w.Header().Set(key, setting.Header.Get(key))
		}
	}

	status := http.StatusOK
	if len(statusCode) > 0 && statusCode[0] > 0 {
		status = statusCode[0]
	}

	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(data)
	return
}

func GetRequestVars(r *http.Request, keys ...string) map[string]string {
	vars := make(map[string]string, len(keys))
	muxVars := mux.Vars(r)
	r.ParseForm()
	for _, key := range keys {
		value := muxVars[key]
		if len(value) == 0 {
			value = r.FormValue(key)
		}
		vars[key] = value
	}
	return vars
}

func ContextGet(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

func ContextSet(r *http.Request, key, val interface{}) *http.Request {
	if val == nil {
		return r
	}
	return r.WithContext(context.WithValue(r.Context(), key, val))
}

func IsTokenRequest(r *http.Request) (bool, string) {
	tokenRequest := r.URL.Path == system.TokenPath
	if tokenRequest {
		grantType := oauth2.GrantType(r.FormValue("grant_type")).String()
		return tokenRequest, grantType
	}
	return tokenRequest, ""
}

func IsProtected(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, system.ProtectedPrefix)
}

func TraceLoger(code string, r *http.Request, data ...map[string]interface{}) *logrus.Entry{
	traceFields := logrus.Fields{
		"trace" : true,
		"code" : code,
	}
	userInfo := GetUserInfo(r)
	if userInfo != nil {
		traceFields["userId"] = userInfo.Id
		traceFields["user"] = userInfo.Code
		traceFields["tenantId"] = userInfo.TenantId
		traceFields["tenant"] = userInfo.TenantCode
	}
	traceFields["requestUrl"] = r.RequestURI

	info := make(map[string]interface{})
	if len(data) > 0 {
		for i := range data {
			dataSet := data[i]
			if len(dataSet) > 0 {
				for key, value := range dataSet {
					info[key] = value
				}
			}
		}
	}
	info["remoteAddr"] = r.RemoteAddr
	if infoJson, err := json.Marshal(info); err == nil {
		traceFields["info"] = string(infoJson)
	}

	return logrus.WithFields(traceFields)
}


package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/fidelfly/gox/authx"
	"github.com/fidelfly/gox/httprxr"
	"github.com/fidelfly/gox/logx"
	"github.com/fidelfly/gox/routex"

	"github.com/fidelfly/fxgos/cmd/service/api/filedb"
)

func FileRoute(router *routex.Router) {
	rr := router.PathPrefix("/file").Subrouter()
	rr.Restricted(true)
	rr.Methods(http.MethodGet).Path("/{id}").HandlerFunc(getFile).Restricted(false)
	rr.Methods(http.MethodGet).HandlerFunc(getFile).Restricted(false)
	rr.Methods(http.MethodPost).HandlerFunc(postFile)
}

func getFile(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "id")
	fileID, _ := strconv.ParseInt(params["id"], 10, 64)
	if fileID == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}
	if resFile, err := filedb.Read(r.Context(), fileID); err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
	} else if resFile == nil {
		httprxr.ResponseJSON(w, http.StatusNotFound, nil)
	} else {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", resFile.Name))
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resFile.Data)
	}
}

func postFile(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUserInfo(r)
	if currentUser == nil {
		httprxr.ResponseJSON(w, http.StatusUnauthorized, httprxr.NewErrorMessage(authx.UnauthorizedErrorCode, "invalid access"))
		return
	}
	key := r.FormValue("key")
	mf, h, err := r.FormFile(key)
	defer func() {
		if err != nil {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
			return
		}
	}()
	if err != nil {
		return
	}
	defer func() {
		logx.CaptureError(mf.Close())
	}()
	data, err := ioutil.ReadAll(mf)
	if err != nil {
		return
	}
	fileId, err := filedb.Save(r.Context(), h.Filename, data)
	if err != nil {
		return
	}
	httprxr.ResponseJSON(w, http.StatusOK, fileId)
}

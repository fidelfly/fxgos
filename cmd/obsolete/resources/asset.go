package resources

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/fidelfly/gox/httprxr"
	"github.com/fidelfly/gox/logx"
	"github.com/fidelfly/gox/pkg/filex"

	"github.com/fidelfly/fxgos/cmd/utilities/system"
)

type ResourceAsset struct {
	ID         int64     `json:"id"`
	Md5        string    `json:"md5"`
	Type       string    `json:"type"`
	Size       int64     `json:"size"`
	Name       string    `json:"name"`
	CreateTime time.Time `json:"create_time"`
}

type AssetService struct {
}

func (as *AssetService) Post(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	mf, h, err := r.FormFile(key)
	defer func() {
		if err != nil {
			logx.Error(err)
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
	asset := new(system.Assets)
	asset.Md5 = filex.CalculateBytesMd5(data)

	find, _ := system.DbEngine.Where("md5 = ? ", asset.Md5).Get(asset)

	if !find {
		asset.Name = h.Filename
		asset.Size = h.Size
		asset.Data = data
		asset.Type = h.Header.Get("Content-Type")

		_, err = system.DbEngine.Insert(asset)
		if err != nil {
			return
		}
	}

	httprxr.ResponseJSON(w, http.StatusOK, ResourceAsset{
		ID:         asset.ID,
		Md5:        asset.Md5,
		Size:       asset.Size,
		Type:       asset.Type,
		Name:       asset.Name,
		CreateTime: asset.CreateTime,
	})

}

func (as *AssetService) Get(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "id")
	id, _ := strconv.ParseInt(params["id"], 10, 64)
	if id == 0 {
		httprxr.ResponseJSON(w, http.StatusBadRequest, httprxr.InvalidParamError("id"))
		return
	}

	asset := system.Assets{
		ID: id,
	}

	find, err := system.DbEngine.Get(&asset)
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	}
	if find {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", asset.Name))
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		logx.CaptureError(w.Write(asset.Data))
	} else {
		httprxr.ResponseJSON(w, http.StatusNotFound, nil)
	}

}

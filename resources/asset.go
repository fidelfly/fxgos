package resources

import (
	"net/http"
	"time"
	"github.com/lyismydg/fxgos/service"
	"github.com/lyismydg/fxgos/system"
	"io/ioutil"
	"github.com/sirupsen/logrus"
	"github.com/lyismydg/fxgos/pkg/file"
	"strconv"
	"fmt"
)

type ResourceAsset struct {
	Id int64 `json:"id"`
	Md5 string `json:"md5"`
	Type string `json:"type"`
	Size int64 `json:"size"`
	Name string `json:"name"`
	CreateTime time.Time `json:"create_time"`
}

type AssetService struct {

}

func (as *AssetService) Post(w http.ResponseWriter, r * http.Request) {
	key := r.FormValue("key")
	mf, h, err := r.FormFile(key)
	defer func() {
		if err != nil {
			logrus.Error(err)
			service.ResponseJSON(w, nil, service.ExceptionError(err), http.StatusInternalServerError)
			return
		}
	}()
	if err != nil {
		return
	}
	defer mf.Close()

	data, err := ioutil.ReadAll(mf)
	if err !=  nil {
		return
	}
	asset := new(system.Assets)
	asset.Md5 = file.GetBytesMd5(data)

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

	assetRes := &ResourceAsset{
		Id: asset.Id,
		Md5: asset.Md5,
		Size: asset.Size,
		Type: asset.Type,
		Name: asset.Name,
		CreateTime: asset.CreateTime,
	}

	service.ResponseJSON(w, nil, assetRes, http.StatusOK)

}

func (as *AssetService) Get(w http.ResponseWriter, r * http.Request) {
	params := service.GetRequestVars(r, "id")
	id, _ := strconv.ParseInt(params["id"], 10, 64)
	if id == 0 {
		service.ResponseJSON(w, nil, service.InvalidParamError("id"), http.StatusBadRequest)
		return
	}

	asset := system.Assets{
		Id: id,
	}

	find, err := system.DbEngine.Get(&asset)
	if err != nil {
		service.ResponseJSON(w, nil, service.ExceptionError(err), http.StatusInternalServerError)
		return
	}
	if find {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", asset.Name))
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write(asset.Data)
	} else {
		service.ResponseJSON(w, nil, nil, http.StatusNotFound)
	}

}

func init()  {
	asset := new(AssetService)
	defineResourceHandlerFunction("post", "/admin/asset", asset.Post)
	defineResourceHandlerFunction("get", "/public/asset/{id}", asset.Get)
}

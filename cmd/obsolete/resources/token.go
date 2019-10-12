package resources

import (
	"net/http"
	"time"

	"github.com/fidelfly/fxgo/httprxr"
	"github.com/fidelfly/fxgo/logx"
)

func logout(w http.ResponseWriter, r *http.Request) {

	userInfo := GetUserInfo(r)

	logx.Infof("%s logout at %s", userInfo.Code, time.Now().Format("2006-01-02 15:04:05"))

	data := httprxr.ResponseData{}
	data["status"] = true

	httprxr.ResponseJSON(w, http.StatusOK, data)
}

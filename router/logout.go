package router

import (
		"github.com/lyismydg/fxgos/service"
	"net/http"
			"time"
)

func logout(w http.ResponseWriter, r *http.Request) {

	userInfo := service.GetUserInfo(r)

	service.TraceLoger("LOGOUT", r).Infof("%s logout at %s", userInfo.Code, time.Now().Format("2006-01-02 15:04:05"))

	data := service.ResponseData{}
	data["status"] = true

	service.ResponseJSON(w, nil, data, http.StatusOK)
}

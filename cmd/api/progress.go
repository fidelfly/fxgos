package api

import (
	"net/http"

	"github.com/fidelfly/fxgo/cachex/mcache"
	"github.com/fidelfly/fxgo/httprxr"
	"github.com/fidelfly/fxgo/logx"
	"github.com/fidelfly/fxgo/pkg/randx"
	"github.com/fidelfly/fxgo/progx"
	"github.com/fidelfly/fxgo/routex"
)

func ProgressRoute(router *routex.Router) {
	rr := router.PathPrefix("/progress").Subrouter()
	rr.Restricted(true)
	rr.Methods(http.MethodGet).Path("/{code}").HandlerFunc(setup)
	rr.Methods(http.MethodGet).HandlerFunc(setup)
}

func setup(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "code")
	code := params["code"]
	wsc := &httprxr.WsConnect{Code: code}
	err := httprxr.SetupWebsocket(wsc, w, r)
	if err != nil {
		httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		return
	}

	progressKey := randx.GetUUID(code)
	SocketCache.Set(progressKey, wsc)
	defer SocketCache.Remove(progressKey)

	_ = wsc.Conn.WriteJSON(map[string]string{"progressKey": progressKey})

	wsc.ListenAndServe()

	logx.Infof("WebSocket %s is Closed", r.RequestURI)
}

var SocketCache = mcache.NewCache(0, 0)

func getProgress(key string, code string) *progx.Progress {
	if conn, ok := SocketCache.Get(key); ok {
		if wsconn, ok := conn.(*httprxr.WsConnect); ok {
			return progx.NewProgress((*httprxr.WsProgressHandler)(wsconn), code)
		}

	}
	return progx.NewProgress(nil, code)
}

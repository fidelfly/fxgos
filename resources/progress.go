package resources

import (
	"github.com/lyismydg/fxgos/service"
	"github.com/lyismydg/fxgos/websocket"
	"net/http"

	"github.com/lyismydg/fxgos/system"

	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type ProgressService struct {
}

func (ws *ProgressService) Setup(w http.ResponseWriter, r *http.Request) {
	params := service.GetRequestVars(r, "code")
	code := params["code"]

	wsc := &websocket.WsConnect{Code: code}

	err := websocket.SetupWebsocket(wsc, w, r)
	if err != nil {
		service.ResponseJSON(w, nil, service.ExceptionError(err), http.StatusInternalServerError)
		return
	}

	progressKey := generateSocketKey(code, r)

	system.SocketCache.Set(progressKey, wsc)
	defer system.SocketCache.Remove(progressKey)

	wsc.Conn.WriteJSON(map[string]string{"progressKey": progressKey})

	wsc.ListenAndServe()

	logrus.Infof("WebSocket %s is Closed", r.RequestURI)
}

func init() {
	ws := &ProgressService{}
	myRouter.Root().Path(service.GetProtectedPath("progress")).HandlerFunc(ws.Setup)
	myRouter.Root().Path(service.GetProtectedPath("progress/{code}")).HandlerFunc(ws.Setup)
}

func generateSocketKey(code string, r *http.Request) string {
	uid, _ := uuid.NewV4()
	key := uuid.Must(uuid.NewV5(uid, code), nil)
	return key.String()
}

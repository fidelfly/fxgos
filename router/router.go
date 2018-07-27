package router

import (
	"github.com/gorilla/mux"
	"github.com/lyismydg/fxgos/auth"
	"github.com/lyismydg/fxgos/resources"
	"net/http"
	"github.com/lyismydg/fxgos/service"
	"gopkg.in/oauth2.v3"
)

func SetupRouter() (router *mux.Router, err error) {
	router = mux.NewRouter()
	err = setupServiceRouter(router)
	if err != nil {
		return
	}
	err = auth.SetupOAuthRouter(router)
	if err != nil {
		return
	}

	resources.SetupRouter(router)

	router.Use(logMiddleware)
	return
}


func setupServiceRouter(router *mux.Router) (err error) {
	router.HandleFunc("/service/tenant", createTeant).Methods("post")
	router.HandleFunc("/admin/logout", logout)
	return
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logWriter := newLoggingResponseWriter(w)
		next.ServeHTTP(logWriter, r)

		logData := make(map[string]interface{})
		logData["status"] = http.StatusText(logWriter.statusCode);
		traceCode := "API_LOG"
		if ok, grantType := service.IsTokenRequest(r); ok {
			switch oauth2.GrantType(grantType) {
			case oauth2.Refreshing:
				traceCode = "TOKEN_REFRESH"
				break
			default:
				traceCode = "TOKEN_GRANT"
				break
			}
		}
		userInfo := service.GetUserInfo(r)
		if userInfo != nil {
			service.TraceLoger(traceCode, r, logData).Infof("%s %s (%s, %s[%v] in %s[%s])", r.Method, r.RequestURI, r.RemoteAddr, userInfo.Code, userInfo.Id, userInfo.TenantName, userInfo.TenantCode)
		} else {
			service.TraceLoger(traceCode, r, logData).Infof("%s %s ( %s )", r.Method, r.RequestURI, r.RemoteAddr)
		}

	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

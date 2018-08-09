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
	router.HandleFunc("/fxgos/logout", logout)
	router.HandleFunc("/fxgos/password", updatePassword).Methods("post")
	router.HandleFunc("/test/{key}", Test)
	router.HandleFunc("/fxgos/test/{key}", Test)
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
		service.TraceLoger(traceCode, r, logData).Infof("%s %s", r.Method, r.RequestURI)
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

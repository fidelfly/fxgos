package main

import (
	"net/http"

	"github.com/fidelfly/fxgo"
	"github.com/fidelfly/fxgo/authx"
	"github.com/fidelfly/fxgo/httprxr"

	"github.com/fidelfly/fxgos/auth"
	_ "github.com/fidelfly/fxgos/example"
	_ "github.com/fidelfly/fxgos/resources"
	"github.com/fidelfly/fxgos/service"
	"github.com/fidelfly/fxgos/system"
)

func setupRouter() (router *fxgo.RootRouter) {
	// enable router audit
	router = fxgo.EnableRouterAudit()

	// setup router authorize route
	client := system.OAuth2.Client[0]
	tokenStore, _ := authx.NewFileTokenStore("./token")
	authServer := authx.SetupPasswordAuthServer(
		&authx.AuthClient{ID: client.ID, Secret: client.Secret, Domain: client.Domain},
		auth.AuthorizationHandler,
		tokenStore,
		authx.WebTokenCfg(),
	)
	fxgo.SetupAuthorizeRoute(system.TokenPath, authServer)

	router.Use(contextMiddleware)

	// setup route table
	fxgo.AttachHookRoute()
	fxgo.SetupProgressRoute(system.GetProtectedPath("/progress"), true)
	return
}

func contextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userKey := httprxr.ContextGet(r, fxgo.ContextUserKey)
		if userKey != nil {
			if key, ok := userKey.(string); ok {
				if userInfo, ok := system.UserCache.EnsureGet(key); ok {
					r = httprxr.ContextSet(r, service.ContextKeys.UserInfo, userInfo)
				}

			}
		}

		next.ServeHTTP(w, r)

	})

}

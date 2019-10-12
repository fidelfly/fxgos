package resources

import (
	"github.com/fidelfly/fxgo/gosrvx"

	"github.com/fidelfly/fxgos/cmd/utilities/system"
)

func init() {
	gosrvx.AddRouterHook(setupRouter)
	setupRouter()
}

func setupRouter() {
	var pRouter = gosrvx.ProtectPrefix(system.ProtectedPrefix)

	// asset
	asset := new(AssetService)
	pRouter.Path("/asset").Methods("post").HandlerFunc(asset.Post)
	gosrvx.Router().Path(system.GetPublicPath("asset/{id}")).Methods("get").HandlerFunc(asset.Get)

	// user
	user := new(UserService)
	pRouter.Path("/user").Methods("get").HandlerFunc(user.Get)
	pRouter.Path("/user").Methods("post").HandlerFunc(user.Post)
	gosrvx.Router().Path(system.GetPublicPath("user")).Methods("post").HandlerFunc(user.Register)
	pRouter.Path("/password").Methods("post").HandlerFunc(user.updatePassword)

	// logout
	pRouter.Path("/logout").Handler(gosrvx.AttachFuncMiddleware(logout, system.TokenKeeper.AuthorizeDisposeMiddleware))
}

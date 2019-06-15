package resources

import (
	"github.com/fidelfly/fxgo"

	"github.com/fidelfly/fxgos/system"
)

func init() {
	fxgo.AddRouterHook(setupRouter)
	setupRouter()
}

func setupRouter() {
	var pRouter = fxgo.ProtectPrefix(system.ProtectedPrefix)

	// asset
	asset := new(AssetService)
	pRouter.Path("/asset").Methods("post").HandlerFunc(asset.Post)
	fxgo.Router().Path(system.GetPublicPath("asset/{id}")).Methods("get").HandlerFunc(asset.Get)

	// user
	user := new(UserService)
	pRouter.Path("/user").Methods("get").HandlerFunc(user.Get)
	pRouter.Path("/user").Methods("post").HandlerFunc(user.Post)
	fxgo.Router().Path(system.GetPublicPath("user")).Methods("post").HandlerFunc(user.Register)
	pRouter.Path("/password").Methods("post").HandlerFunc(user.updatePassword)

	// logout
	pRouter.Path("/logout").Handler(fxgo.AttachFuncMiddleware(logout, system.TokenKeeper.AuthorizeDisposeMiddleware))
}

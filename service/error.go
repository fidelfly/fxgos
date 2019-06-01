package service

import (
	"github.com/fidelfly/fxgo/authx"
	"github.com/fidelfly/fxgo/httprxr"
)

var UnauthorizedError = httprxr.NewErrorMessage(authx.UnauthorizedErrorCode, "UNAUTHORIZED")
var TokenExpiredError = httprxr.NewErrorMessage(authx.TokenExpiredErrorCode, "Token is expired")
var NotSupportError = httprxr.NewErrorMessage("not_support", "Function is not support!")

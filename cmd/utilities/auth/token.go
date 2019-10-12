package auth

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"strconv"

	"github.com/fidelfly/fxgo/authx"
	"github.com/fidelfly/fxgo/gosrvx"

	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/user/res"
)

var TokenIssuer *gosrvx.TokenIssuer

func setupAuthServer(clients ...authx.AuthClient) (server *authx.Server, err error) {
	tokenStore, err := authx.NewFileTokenStore("./token")
	if err != nil {
		return
	}
	server = authx.NewOAuthServer()
	server.SetTokenStorage(tokenStore)
	cps := make([]authx.ClientInfo, len(clients))
	for index, client := range clients {
		cps[index] = &client
	}
	server.SetClients(cps...)
	server.SetPasswordAuthorizationHandler(passwordHandler)
	server.SetExtensionFieldsHandler(authx.NewTokenExtension(authx.UserExtension))

	return
}

func SetupTokenIssuer(tokenPath string, clients ...authx.AuthClient) (*gosrvx.TokenIssuer, error) {
	server, err := setupAuthServer(clients...)
	if err != nil {
		return nil, err
	}
	TokenIssuer = gosrvx.NewTokenIssuer(server, tokenPath)
	return TokenIssuer, nil
}

func passwordHandler(username, password string) (userID string, err error) {
	user := res.User{}
	encodedPwd := EncodePassword(username, password)
	ok, err := db.Read(&user, db.Where("code = ? and password = ?", username, encodedPwd))
	if ok {
		userID = strconv.FormatInt(user.Id, 10)
	}
	return
}

func EncodePassword(code string, pasword string) string {
	plainPwd := fmt.Sprintf("%s:%s", code, pasword)
	data := sha256.Sum256([]byte(plainPwd))
	encodedPwd := md5.Sum(data[:])
	return fmt.Sprintf("%x", encodedPwd)
}

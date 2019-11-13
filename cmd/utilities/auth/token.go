package auth

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	"github.com/fidelfly/gox/authx"
	"github.com/fidelfly/gox/gosrvx"

	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/gostool/db"
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
	for index := range clients {
		cps[index] = &clients[index]
	}
	server.SetClients(cps...)
	server.SetPasswordAuthorizationHandler(passwordHandler)
	server.SetPasswordTokenCfg(2*time.Hour, 3*time.Hour, true)
	server.SetRefreshTokenCfg(2*time.Hour, 3*time.Hour, true, true, true, true)
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

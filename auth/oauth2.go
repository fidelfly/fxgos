package auth

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/lyismydg/fxgos/service"
	"github.com/lyismydg/fxgos/system"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

type AuthServerConfig struct {
	config     system.OAuth2Properties
	authServer *server.Server
}

func (asc *AuthServerConfig) Setup() (err error) {
	authManager := manage.NewDefaultManager()
	authManager.MustTokenStorage(store.NewMemoryTokenStore())
	authManager.SetPasswordTokenCfg(&manage.Config{AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 3, IsGenerateRefresh: true})
	authManager.SetRefreshTokenCfg(&manage.RefreshingConfig{IsGenerateRefresh: true, IsRemoveAccess: true, IsRemoveRefreshing: true, IsResetRefreshTime: true})
	clientStore := store.NewClientStore()

	for _, client := range asc.config.Client {
		clientStore.Set(client.Id, &models.Client{
			ID:     client.Id,
			Secret: client.Secret,
			Domain: client.Domain,
		})
	}

	authManager.MapClientStorage(clientStore)

	asc.authServer = server.NewDefaultServer(authManager)

	asc.authServer.PasswordAuthorizationHandler = authorizationHandler
	asc.authServer.ExtensionFieldsHandler = tokenFields
	return
}

func (asc *AuthServerConfig) authValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenInfo oauth2.TokenInfo
		if service.IsProtected(r) {

			if ti, ok := asc.validateToken(w, r); !ok {
				return
			} else if ti != nil {
				if userInfo, ok := system.UserCache.EnsureGet(ti.GetUserID()); ok {
					r = service.ContextSet(r, service.ContextKeys.UserInfo, userInfo)
				}
				tokenInfo = ti
			}

		}

		next.ServeHTTP(w, r)

		if r.URL.Path == service.GetProtectedPath("logout") && tokenInfo != nil {
			asc.authServer.Manager.RemoveAccessToken(tokenInfo.GetAccess())
			asc.authServer.Manager.RemoveRefreshToken(tokenInfo.GetRefresh())
		}
	})
}

func (asc *AuthServerConfig) validateToken(w http.ResponseWriter, r *http.Request) (ti oauth2.TokenInfo, status bool) {
	tokenRequest, grantType := service.IsTokenRequest(r)

	if tokenRequest && oauth2.GrantType(grantType) != oauth2.Refreshing {
		status = true
		return
	}

	ti, err := asc.authServer.ValidationBearerToken(r)
	if err != nil {
		service.ResponseJSON(w, nil, service.UnauthorizedError, http.StatusUnauthorized)
		status = false
		return
	}

	if ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).After(time.Now()) {
		status = true
		return
	}

	if ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(time.Now()) {
		service.ResponseJSON(w, nil, service.UnauthorizedError, http.StatusUnauthorized)
		status = false
		return
	}

	if tokenRequest && oauth2.GrantType(grantType) == oauth2.Refreshing {
		status = true
		return
	}

	service.ResponseJSON(w, nil, service.TokenExpired, http.StatusUnauthorized)
	status = false
	return
}

func (asc *AuthServerConfig) handlerTokenRequest(w http.ResponseWriter, r *http.Request) {
	asc.authServer.HandleTokenRequest(w, r)
}

func EncodePassword(code string, pasword string) string {
	plainPwd := fmt.Sprintf("%s:%s", code, pasword)
	data := sha256.Sum256([]byte(plainPwd))
	encodedPwd := md5.Sum(data[:])
	return fmt.Sprintf("%x", encodedPwd)
}

func authorizationHandler(username, password string) (userID string, err error) {
	user := system.User{}
	encodedPwd := EncodePassword(username, password)
	ok, err := system.DbEngine.Where("code = ? and password = ?", username, encodedPwd).Get(&user)
	if ok {
		userID = strconv.FormatInt(user.Id, 10)
	}
	return
}

func tokenFields(ti oauth2.TokenInfo) (fieldsValue map[string]interface{}) {
	data := make(map[string]interface{}, 1)
	data["user_id"], _ = strconv.Atoi(ti.GetUserID())
	return data
}

func SetupOAuthRouter(router *mux.Router) (err error) {
	authServerConfig := AuthServerConfig{config: *system.OAuth2}
	err = authServerConfig.Setup()
	if err != nil {
		return
	}

	router.HandleFunc(system.TokenPath, authServerConfig.handlerTokenRequest)
	router.Use(authServerConfig.authValidateMiddleware)
	return
}

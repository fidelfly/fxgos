package auth

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"strconv"

	"github.com/fidelfly/fxgos/system"
)

func EncodePassword(code string, pasword string) string {
	plainPwd := fmt.Sprintf("%s:%s", code, pasword)
	data := sha256.Sum256([]byte(plainPwd))
	encodedPwd := md5.Sum(data[:])
	return fmt.Sprintf("%x", encodedPwd)
}

func AuthorizationHandler(username, password string) (userID string, err error) {
	user := system.User{}
	encodedPwd := EncodePassword(username, password)
	ok, err := system.DbEngine.Where("code = ? and password = ?", username, encodedPwd).Get(&user)
	if ok {
		userID = strconv.FormatInt(user.ID, 10)
	}
	return
}

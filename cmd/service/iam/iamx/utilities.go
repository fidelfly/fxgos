package iamx

import (
	"fmt"
	"strconv"
	"strings"
)

//export
func EncodeObject(resType string, code string, keys ...interface{}) string {
	return fmt.Sprintf("%s.%s%s", resType, code, fmt.Sprintf(strings.Repeat("_%v", len(keys)), keys...))
}

//export
func DecodeObject(sub string) (resType string, code string, keys []string) {
	index := strings.LastIndex(sub, ".")
	if index <= 0 {
		return
	}
	resType = sub[:index]
	resKey := strings.Split(sub[index+1:], "_")
	code = resKey[0]
	if len(resKey) > 1 {
		keys = resKey[1:]
	}
	return
}

//export
func EncodeUserSubject(userId int64) string {
	return fmt.Sprintf("user_%d", userId)
}

//export
func EncodeRoleSubject(roleId int64) string {
	return fmt.Sprintf("role_%d", roleId)
}

//export
func DecodeUserSubject(sub string) int64 {
	subs := strings.Split(sub, "_")
	if len(subs) == 2 {
		if id, err := strconv.ParseInt(subs[1], 10, 64); err == nil {
			return id
		}
	}
	return 0
}

//export
func DecodeRoleSubject(sub string) int64 {
	subs := strings.Split(sub, "_")
	if len(subs) == 2 {
		if id, err := strconv.ParseInt(subs[1], 10, 64); err == nil {
			return id
		}
	}
	return 0
}

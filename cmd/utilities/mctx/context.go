package mctx

import (
	"context"
	"reflect"
	"strconv"

	"github.com/fidelfly/gox/gosrvx"
)

func GetUserId(ctx context.Context) int64 {
	userKey := ctx.Value(gosrvx.ContextUserKey)
	if userKey == nil {
		return 0
	}
	if key, ok := userKey.(string); ok {
		id, err := strconv.ParseInt(key, 10, 64)
		if err == nil {
			return id
		}
	}
	return 0
}

func FillUserInfo(ctx context.Context, target interface{}) bool {
	userId := GetUserId(ctx)
	if userId > 0 {
		setUserInfo(target, userId)
		return true
	}
	return false
}

func setUserInfo(target interface{}, userId int64) {
	v := reflect.ValueOf(target)
	if v.IsValid() == false {
		return
	}

	if v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}

	s := v.Elem()
	if s.Kind() == reflect.Struct {
		f := s.FieldByName("CreateUser")
		if f.IsValid() {
			if f.CanSet() && f.Kind() == reflect.Int64 {
				if f.Int() == 0 {
					f.SetInt(userId)
				}
			}

		}

		f = s.FieldByName("UpdateUser")
		if f.IsValid() {
			if f.CanSet() && f.Kind() == reflect.Int64 {
				f.SetInt(userId)
			}

		}
	}
}

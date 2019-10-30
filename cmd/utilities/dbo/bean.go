package dbo

import (
	"context"
	"reflect"

	"github.com/fidelfly/gox/pkg/reflectx"
	"github.com/fidelfly/gox/pkg/strx"

	"github.com/fidelfly/fxgos/cmd/utilities/mctx"
)

type BeanOption func(interface{})

func CreateUserOption(ctx context.Context, fields ...string) BeanOption {
	userId := mctx.GetUserId(ctx)
	if userId == 0 {
		return nil
	}
	if strx.IndexOfSlice(fields, "CreateUser") < 0 {
		fields = append(fields, "CreateUser")
	}
	if strx.IndexOfSlice(fields, "UpdateUser") < 0 {
		fields = append(fields, "UpdateUser")
	}
	return func(target interface{}) {
		pairs := make([]reflectx.FV, len(fields))
		for i, f := range fields {
			pairs[i] = reflectx.FV{
				Field: f,
				Value: userId,
			}
		}
		reflectx.SetField(target, pairs...)
	}
}

func ApplyBeanOption(target interface{}, option ...BeanOption) interface{} {
	sliceValue := reflect.Indirect(reflect.ValueOf(target))
	if sliceValue.Kind() == reflect.Slice {
		size := sliceValue.Len()
		if size > 0 {
			for i := 0; i < size; i++ {
				sv := sliceValue.Index(i)
				if sv.Kind() == reflect.Ptr {
					_applyBeanOption(sv.Interface(), option...)
				} else if sv.CanAddr() {
					_applyBeanOption(sv.Addr().Interface(), option...)
				}

			}
		}
	} else if sliceValue.Kind() == reflect.Struct {
		_applyBeanOption(target, option...)
	}
	return target
}

func _applyBeanOption(target interface{}, opts ...BeanOption) {
	for _, opt := range opts {
		opt(target)
	}
}

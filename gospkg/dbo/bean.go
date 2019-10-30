package dbo

import (
	"reflect"

	"github.com/fidelfly/gox/pkg/reflectx"
)

type BeanOption interface {
	Apply(interface{})
}

type FuncBeanOption func(interface{})

func (fbo FuncBeanOption) Apply(t interface{}) {
	fbo(t)
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
		opt.Apply(target)
	}
}

func Assignment(s interface{}) BeanOption {
	return FuncBeanOption(func(t interface{}) {
		if t != s {
			reflectx.CopyAllFields(t, s)
		}
	})
}

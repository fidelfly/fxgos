package mdbo

import "github.com/fidelfly/gox/pkg/reflectx"

func getId(target interface{}) int64 {
	if v := reflectx.GetField(target, "Id"); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

func getCode(target interface{}) string {
	if v := reflectx.GetField(target, "Code"); v != nil {
		if code, ok := v.(string); ok {
			return code
		}
	}
	return ""
}

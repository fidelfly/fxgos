package service

import (
	"fmt"
	"github.com/lyismydg/fxgos/system"
)

func InvalidParamError(paramName string, data... map[string] interface{}) ResponseError {
	return ResponseError{
		Code:"INVALID_PARAMS",
		Message: fmt.Sprintf("Param:%s is invalid!", paramName),
		Data: combinedData(data...),
	}
}

func ExceptionError(err error, data... map[string]interface{}) ResponseError {
	return ResponseError{
		Code: "EXCEPTION",
		Message: err.Error(),
		Data: combinedData(data...),
	}
}

func ResourceLockedError(action *system.LockAction) ResponseError  {
	var data map[string] interface{}
	if action != nil {
		data = make(map[string] interface{})
		user := system.User{
			Id: action.UserId,
		}
		_, err := system.DbEngine.Get(&user)
		if err != nil {
			data["user"] = action.UserId
		} else {
			data["user"] = user.Name
		}
		data["action"] = action.Code
	}

	return NewResponseError("RESOURCE_LOCKED", "Resource is locked by someone. Please try again later.", data)
}

func NewResponseError(code string, message string, data... map[string]interface{}) ResponseError {
	return ResponseError{
		Code: code,
		Message: message,
		Data: combinedData(data...),
	}
}

func combinedData(data... map[string]interface{}) map[string] interface{}  {
	if len(data) == 0 {
		return nil
	} else if len(data) == 1 {
		return data[0]
	}

	mapData := make(map[string]interface{})

	for _, m := range data {
		for key, value := range m {
			mapData[key] = value
		}
	}

	return mapData

}

var UnauthorizedError = NewResponseError("UNAUTHORIZED", "Unauthorized Action!")
var TokenExpired = NewResponseError("TOKENEXPIRED", "Token is expired!")
var NotSupport = NewResponseError("NOT_SUPPORT", "Function is not support!")
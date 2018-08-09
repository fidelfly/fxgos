package service

import "fmt"

func InvalidParamError(paramName string) ResponseError {
	return ResponseError{
		"INVALID_PARAMS",
		fmt.Sprintf("Param:%s is invalid!", paramName),
	}
}

func ExceptionError(err error) ResponseError {
	return ResponseError{
		"EXCEPTION",
		err.Error(),
	}
}

func NewResponseError(code string, message string) ResponseError {
	return ResponseError{
		code,
		message,
	}
}

var UnauthorizedError = ResponseError{"UNAUTHORIZED", "Unauthorized Action!"}
var TokenExpired = ResponseError{"TOKENEXPIRED", "Token is expired!"}

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

var UnauthorizedError = ResponseError{"UNAUTHORIZED", "Unauthorized Action!"}
var TokenExpired = ResponseError{"TOKENEXPIRED", "Token is expired!"}

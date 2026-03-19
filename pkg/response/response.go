package response

import "net/http"

type CommonResult struct {
	Code    uint
	Message string
	Data    any
}

func commonResult(code uint, message string, data interface{}) *CommonResult {
	return &CommonResult{Code: code, Message: message, Data: data}
}

func Success(data interface{}) *CommonResult {
	return commonResult(http.StatusOK, "请求成功", data)
}

func Fail() *CommonResult {
	return commonResult(http.StatusInternalServerError, "系统异常，请稍后重试", nil)
}

func FailWithMsg(message string) *CommonResult {
	return commonResult(http.StatusInternalServerError, message, nil)
}

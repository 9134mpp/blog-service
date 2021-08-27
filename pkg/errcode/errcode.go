package errcode

import (
	"fmt"
	"net/http"
)

var codes = map[int]string{}

type Error struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Details []string `json:"details"` // 详细信息
}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
	}
	codes[code] = msg
	return &Error{Code: code, Msg: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("错误码：%d, 错误信息:：%s", e.GetCode(), e.GetMsg())
}

func (e *Error) GetMsg() string {
	return e.Msg
}

func (e *Error) GetCode() int {
	return e.Code
}

func (e *Error) MsgF(args []interface{}) string {
	return fmt.Sprintf(e.Msg, args...)
}

func (e *Error) GetDetails() []string {
	return e.Details
}

func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	newError.Details = []string{}

	for _, d := range details {
		newError.Details = append(newError.Details, d)
	}

	return &newError
}

func (e *Error) StatusCode() int {
	switch e.GetCode() {
	case Success.GetCode():
		return http.StatusOK
	case ServerError.GetCode():
		return http.StatusInternalServerError
	case InvalidParams.GetCode():
		return http.StatusBadRequest
	case UnauthorizedAuthNotExist.GetCode():
		fallthrough
	case UnauthorizedTokenError.GetCode():
		fallthrough
	case UnauthorizedTokenGenerate.GetCode():
		fallthrough
	case UnauthorizedTokenTimeout.GetCode():
		return http.StatusUnauthorized
	case TooManyRequests.GetCode():
		return http.StatusTooManyRequests
	case NotFound.GetCode():
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

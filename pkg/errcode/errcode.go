package errcode

import (
	"fmt"
	"github.com/guidoxie/knife/pkg/log"
	"net/http"
)

type Error struct {
	code    int      `json:"code"`
	msg     string   `json:"msg"`
	details []string `json:"details"`
}

var codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		log.Panicf("错误码 %d 已经存在，请更换一个", code)
	}
	codes[code] = msg
	return &Error{code: code, msg: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("错误码：%d，错误信息：%s", e.Code(), e.Msg())
}

func (e *Error) Code() int {
	return e.code
}
func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) MsgF(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}

func (e *Error) Details() []string {
	return e.details
}
func (e Error) WithDetails(details ...string) *Error {
	e.details = []string{}
	for _, d := range details {
		e.details = append(e.details, d)
	}
	return &e
}

func (e *Error) StatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case ServerError.Code():
		return http.StatusInternalServerError
	case InvalidParams.Code():
		return http.StatusBadRequest
	case UnauthorizedAuthNotExist.Code():
		fallthrough
	case UnauthorizedAuthTokenError.Code():
		fallthrough
	case UnauthorizedAuthTokenGenerate.Code():
		fallthrough
	case UnauthorizedAuthTokenTimeOut.Code():
		return http.StatusUnauthorized
	case TooManyRequests.Code():
		return http.StatusTooManyRequests
	case NotFound.Code():
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
